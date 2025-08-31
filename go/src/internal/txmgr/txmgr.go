package txmgr

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Manager struct {
	DB *sql.DB
	// 선택: 관측/로깅/트레이싱 훅
	OnBegin    func(ctx context.Context, opts *sql.TxOptions)
	OnCommit   func(ctx context.Context, dur time.Duration, err error)
	OnRollback func(ctx context.Context, err error)
	// 선택: DB 에러 특성(재시도 판단)에 쓰일 헬퍼
	Flavor DBFlavor
}

type Opts struct {
	Isolation  sql.IsolationLevel
	ReadOnly   bool
	MaxRetries int                           // 0=재시도 없음
	Backoff    func(retry int) time.Duration // nil이면 고정 50ms 등
}

type DBFlavor interface {
	// DB/드라이버 별 Deadlock, Serialization failure 등을 판별
	IsRetryableTxError(error) bool
}

// 기본 백오프
func defaultBackoff(_ int) time.Duration { return 50 * time.Millisecond }

func (m *Manager) WithinTx(
	ctx context.Context,
	opts Opts,
	fn func(ctx context.Context, uow *UoW) error,
) error {
	if opts.Backoff == nil {
		opts.Backoff = defaultBackoff
	}
	var lastErr error
	for retry := 0; retry <= opts.MaxRetries; retry++ {
		start := time.Now()
		tx, err := m.DB.BeginTx(ctx, &sql.TxOptions{
			Isolation: opts.Isolation,
			ReadOnly:  opts.ReadOnly,
		})
		if err != nil {
			return err
		}
		if m.OnBegin != nil {
			m.OnBegin(ctx, &sql.TxOptions{Isolation: opts.Isolation, ReadOnly: opts.ReadOnly})
		}

		// UoW: 같은 tx Runner로 모두 묶음
		uow := NewUoW(tx)

		// 비즈니스 실행
		runErr := fn(ctx, uow)

		if runErr != nil {
			_ = tx.Rollback()
			if m.OnRollback != nil {
				m.OnRollback(ctx, runErr)
			}
			// 재시도 판단
			if m.Flavor != nil && m.Flavor.IsRetryableTxError(runErr) && retry < opts.MaxRetries {
				time.Sleep(opts.Backoff(retry))
				lastErr = runErr
				continue
			}
			return runErr
		}

		commitErr := tx.Commit()
		if m.OnCommit != nil {
			m.OnCommit(ctx, time.Since(start), commitErr)
		}
		if commitErr != nil {
			// 커밋 시 실패도 재시도 고려(드물지만)
			if m.Flavor != nil && m.Flavor.IsRetryableTxError(commitErr) && retry < opts.MaxRetries {
				time.Sleep(opts.Backoff(retry))
				lastErr = commitErr
				continue
			}
			return commitErr
		}
		return nil
	}
	// 여기 도달하면 재시도 모두 실패
	if lastErr == nil {
		lastErr = errors.New("transaction failed after retries")
	}
	return lastErr
}
