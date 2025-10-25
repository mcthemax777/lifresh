package worker

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"lifresh/internal/repository/model"
	"log"
	"time"
)

type OutboxWorker struct {
	DB    *gorm.DB
	Redis *redis.Client
	Ctx   context.Context
}

func NewOutboxWorker(db *gorm.DB, rdb *redis.Client) *OutboxWorker {
	return &OutboxWorker{
		DB:    db,
		Redis: rdb,
		Ctx:   context.Background(),
	}
}

func (w *OutboxWorker) Run(interval time.Duration) {
	log.Println("[OutboxWorker] started...")
	for {
		if err := w.processBatch(); err != nil {
			log.Printf("[OutboxWorker] error: %v\n", err)
		}
		time.Sleep(interval)
	}
}

func (w *OutboxWorker) processBatch() error {
	var outboxes []model.Outbox

	// 1️⃣ 미처리 이벤트 가져오기
	if err := w.DB.Where("processed = ?", false).
		Order("id asc").Limit(50).Find(&outboxes).Error; err != nil {
		return err
	}
	if len(outboxes) == 0 {
		return nil
	}

	for _, evt := range outboxes {
		payloadJSON, err := json.Marshal(evt.Payload)
		if err != nil {
			log.Printf("[OutboxWorker] JSON marshal error: %v", err)
			continue
		}

		// 2️⃣ Redis Stream에 발행
		xid, err := w.Redis.XAdd(&redis.XAddArgs{
			Stream: evt.Stream,
			Values: map[string]interface{}{
				"event_id": evt.EventID,
				"payload":  string(payloadJSON),
			},
		}).Result()

		if err != nil {
			log.Printf("[OutboxWorker] Redis publish failed: %v", err)
			continue
		}

		log.Printf("[OutboxWorker] Published event %s -> stream=%s xid=%s",
			evt.EventID, evt.Stream, xid)

		// 3️⃣ DB 업데이트
		if err := w.DB.Model(&model.Outbox{}).
			Where("id = ?", evt.ID).
			Update("processed", true).Error; err != nil {
			log.Printf("[OutboxWorker] mark processed failed: %v", err)
		}
	}
	return nil
}
