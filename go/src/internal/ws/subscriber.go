// internal/ws_gateway/subscriber.go
package ws

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
)

type StreamSubscriber struct {
	Rdb         *redis.Client
	Broadcaster *Broadcaster
}

func NewStreamSubscriber(rdb *redis.Client, b *Broadcaster) *StreamSubscriber {
	return &StreamSubscriber{Rdb: rdb, Broadcaster: b}
}

func (s *StreamSubscriber) SubscribeStream(stream string) {
	group := "ws_group"
	consumer := "ws_consumer"

	s.Rdb.XGroupCreateMkStream(stream, group, "$")

	log.Printf("[Redis] Subscribing to stream: %s", stream)
	for {
		msgs, err := s.Rdb.XReadGroup(&redis.XReadGroupArgs{
			Group:    group,
			Consumer: consumer,
			Streams:  []string{stream, ">"},
			Block:    0,
			Count:    10,
		}).Result()

		if err != nil {
			log.Printf("[Redis] XReadGroup error: %v", err)
			continue
		}

		for _, m := range msgs {
			for _, entry := range m.Messages {
				raw, _ := json.Marshal(entry.Values)
				var sm StreamMessage
				if err := json.Unmarshal(raw, &sm); err != nil {
					continue
				}
				log.Printf("[Redis] Event: %s %s", sm.Entity, sm.Type)

				if sm.UserID != "" {
					s.Broadcaster.BroadcastToUser(sm.UserID, raw)
				}
				if sm.GroupID != "" {
					s.Broadcaster.BroadcastToGroup(sm.GroupID, raw)
				}
			}
		}
	}
}
