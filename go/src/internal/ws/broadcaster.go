package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Broadcaster struct {
	userClients  map[string]map[*websocket.Conn]bool // userId → conn set
	groupClients map[string]map[*websocket.Conn]bool // groupId → conn set
	lock         sync.Mutex
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		userClients:  make(map[string]map[*websocket.Conn]bool),
		groupClients: make(map[string]map[*websocket.Conn]bool),
	}
}

func (b *Broadcaster) AddUserClient(userID string, conn *websocket.Conn) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.userClients[userID]; !ok {
		b.userClients[userID] = make(map[*websocket.Conn]bool)
	}
	b.userClients[userID][conn] = true
	log.Printf("[WS] User %s connected (%d conns)", userID, len(b.userClients[userID]))
}

func (b *Broadcaster) AddGroupClient(groupID string, conn *websocket.Conn) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.groupClients[groupID]; !ok {
		b.groupClients[groupID] = make(map[*websocket.Conn]bool)
	}
	b.groupClients[groupID][conn] = true
	log.Printf("[WS] Joined group %s (%d conns)", groupID, len(b.groupClients[groupID]))
}

func (b *Broadcaster) RemoveUserClient(userID string, conn *websocket.Conn) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if conns, ok := b.userClients[userID]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(b.userClients, userID)
		}
	}
	conn.Close()
}

func (b *Broadcaster) RemoveGroupClient(groupID string, conn *websocket.Conn) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if conns, ok := b.groupClients[groupID]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(b.groupClients, groupID)
		}
	}
	conn.Close()
}

func (b *Broadcaster) BroadcastToUser(userID string, msg []byte) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for conn := range b.userClients[userID] {
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func (b *Broadcaster) BroadcastToGroup(groupID string, msg []byte) {
	b.lock.Lock()
	defer b.lock.Unlock()
	for conn := range b.groupClients[groupID] {
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}
