package ws

type StreamMessage struct {
	Type      string `json:"type"`      // e.g. "update", "delete", "create"
	Entity    string `json:"entity"`    // e.g. "plan", "record"
	EntityID  string `json:"entityId"`  // e.g. "12345"
	UserID    string `json:"userId"`    // 송신 대상 (optional)
	GroupID   string `json:"groupId"`   // 협업방 또는 plan ID
	Payload   any    `json:"payload"`   // 데이터 본문
	Timestamp int64  `json:"timestamp"` // Unix timestamp
}
