package main

import (
	"github.com/go-redis/redis"
	"log"
	"net/http"

	"lifresh/internal/ws"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	broadcaster := ws.NewBroadcaster()
	subscriber := ws.NewStreamSubscriber(rdb, broadcaster)

	// Redis 구독 실행 (비동기)
	go subscriber.SubscribeStream("change_log") // Outbox 에서 발행하는 Stream 이름

	// WebSocket 서버 시작
	server := ws.NewWebSocketServer(broadcaster)
	http.HandleFunc("/ws", server.HandleConnection)

	log.Println("[WebSocket] Server started at :8081")
	_ = http.ListenAndServe(":8081", nil)
}
