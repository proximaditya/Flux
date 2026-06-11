package main

// The message format expected from the Frontend via WebSocket
type WSMessage struct {
	Action   string `json:"action"`             // "invite_friend", "join_room", "get_room_info"
	TargetID string `json:"target_id,omitempty"` // The unique hash of the friend
	RoomID   string `json:"room_id,omitempty"`
	Role     string `json:"role,omitempty"`      // "speaker" or "spectator"
}

// The message format Go sends back to the Frontend
type WSPayload struct {
	Type    string      `json:"type"`    // "invite_received", "room_info", "token_generated"
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}