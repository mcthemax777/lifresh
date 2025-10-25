package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	Broadcaster *Broadcaster
	Upgrader    websocket.Upgrader
}

func NewWebSocketServer(b *Broadcaster) *WebSocketServer {
	return &WebSocketServer{
		Broadcaster: b,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (s *WebSocketServer) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[WS] Upgrade error:", err)
		return
	}
	defer conn.Close()

	var userID, groupID string

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WS] Connection closed for user=%s group=%s", userID, groupID)

			// 연결 종료 시 등록 해제
			if userID != "" {
				s.Broadcaster.RemoveUserClient(userID, conn)
			}
			if groupID != "" {
				s.Broadcaster.RemoveGroupClient(groupID, conn)
			}
			break
		}

		var init map[string]string
		json.Unmarshal(msg, &init)

		if uid, ok := init["userId"]; ok {
			userID = uid
			s.Broadcaster.AddUserClient(uid, conn)
		}
		if gid, ok := init["groupId"]; ok {
			groupID = gid
			s.Broadcaster.AddGroupClient(gid, conn)
		}
	}
}
