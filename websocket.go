package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// Local memory map of connected WebSockets (To push messages to specific users)
var clients = make(map[string]*websocket.Conn)
var clientsMutex sync.RWMutex

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// In production, verify auth token. For MVP, we pass user_id in URL.
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS Upgrade Error:", err)
		return
	}
	defer conn.Close()

	// Register Client locally & in Redis
	clientsMutex.Lock()
	clients[userID] = conn
	clientsMutex.Unlock()

	SetUserOnline(userID)
	log.Printf("🟢 User %s is online", userID)

	// Clean up when they disconnect
	defer func() {
		clientsMutex.Lock()
		delete(clients, userID)
		clientsMutex.Unlock()
		SetUserOffline(userID)
		log.Printf("🔴 User %s went offline", userID)
	}()

	// Listen for messages from Frontend
	for {
		var msg WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		switch msg.Action {

		case "invite_friend":
			// User wants to invite a friend
			targetID := msg.TargetID
			roomID := msg.RoomID

			// Check if target is connected to this specific server node
			clientsMutex.RLock()
			targetConn, exists := clients[targetID]
			clientsMutex.RUnlock()

			if exists {
				// Send invite payload to the friend
				invitePayload := WSPayload{
					Type:    "invite_received",
					Message: fmt.Sprintf("User %s invited you to room %s", userID, roomID),
					Data:    map[string]string{"room_id": roomID, "from": userID},
				}
				targetConn.WriteJSON(invitePayload)
				conn.WriteJSON(WSPayload{Type: "info", Message: "Invite sent!"})
			} else {
				// (In a multi-server setup, you would use Redis Pub/Sub here to find them on another server)
				conn.WriteJSON(WSPayload{Type: "error", Message: "User is offline or on another server"})
			}

		case "get_room_info":
			// User wants to see how many people are in the room before joining
			count, _ := rdb.SCard(ctx, "room:"+msg.RoomID+":members").Result()
			conn.WriteJSON(WSPayload{
				Type: "room_info",
				Data: map[string]interface{}{"room_id": msg.RoomID, "member_count": count},
			})

		case "join_room":
			// User accepted invite or clicked join
			roomID := msg.RoomID
			role := msg.Role // "speaker" or "spectator"

			// 1. Generate LiveKit Token (Mocked function, use previous LiveKit logic here)
			token := "mock_livekit_token_for_" + role

			// 2. Update Redis Room State
			totalMembers := JoinRoomState(roomID, userID)

			// 3. Send Success back to client
			conn.WriteJSON(WSPayload{
				Type: "join_success",
				Data: map[string]interface{}{
					"token":         token,
					"role":          role,
					"total_members": totalMembers,
				},
			})
		}
	}
}
