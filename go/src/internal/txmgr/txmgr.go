package txmgr

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Manager struct {
	DB *gorm.DB
	// 선택: 관측/로깅/트레이싱 훅
	OnBegin    func(opts *sql.TxOptions)
	OnCommit   func(dur time.Duration, err error)
	OnRollback func(err error)
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

type MySQLFlavor struct{}

func (MySQLFlavor) IsRetryableTxError(err error) bool {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	switch me.Number {
	case 1213: // Deadlock found
		return true
	case 1205: // Lock wait timeout
		return true
	default:
		return false
	}
}

// 기본 백오프
func defaultBackoff(_ int) time.Duration { return 50 * time.Millisecond }

func (m *Manager) WithinTx(
	opts Opts,
	fn func(tx *gorm.DB) error,
) error {
	if opts.Backoff == nil {
		opts.Backoff = defaultBackoff
	}
	var lastErr error
	for retry := 0; retry <= opts.MaxRetries; retry++ {
		start := time.Now()
		tx := m.DB.Begin(&sql.TxOptions{
			Isolation: opts.Isolation,
			ReadOnly:  opts.ReadOnly,
		})
		if tx.Error != nil {
			return tx.Error
		}
		if m.OnBegin != nil {
			m.OnBegin(&sql.TxOptions{Isolation: opts.Isolation, ReadOnly: opts.ReadOnly})
		}

		// 비즈니스 실행
		runErr := fn(tx)

		if runErr != nil {
			_ = tx.Rollback()
			if m.OnRollback != nil {
				m.OnRollback(runErr)
			}
			// 재시도 판단
			if m.Flavor != nil && m.Flavor.IsRetryableTxError(runErr) && retry < opts.MaxRetries {
				time.Sleep(opts.Backoff(retry))
				lastErr = runErr
				continue
			}
			return runErr
		}

		tx = tx.Commit()
		if m.OnCommit != nil {
			m.OnCommit(time.Since(start), tx.Error)
		}
		if tx.Error != nil {
			// 커밋 시 실패도 재시도 고려(드물지만)
			if m.Flavor != nil && m.Flavor.IsRetryableTxError(tx.Error) && retry < opts.MaxRetries {
				time.Sleep(opts.Backoff(retry))
				lastErr = tx.Error
				continue
			}
			return tx.Error
		}
		return nil
	}
	// 여기 도달하면 재시도 모두 실패
	if lastErr == nil {
		lastErr = errors.New("transaction failed after retries")
	}
	return lastErr
}
