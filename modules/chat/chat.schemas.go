package chat

import "time"

// ==================== DB Models ====================

// RoomSchema is the MongoDB document for a chat room.
// Think of it as a group chat — it has a name, a creator, and a list of members.
type RoomSchema struct {
	RoomID    string    `bson:"room_id" json:"room_id"`
	Name      string    `bson:"name" json:"name"`
	CreatedBy string    `bson:"created_by" json:"created_by"`
	Members   []string  `bson:"members" json:"members"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// MessageSchema is the MongoDB document for a single chat message.
// Every message belongs to a room and has a sender.
type MessageSchema struct {
	MessageID string    `bson:"message_id" json:"message_id"`
	RoomID    string    `bson:"room_id" json:"room_id"`
	SenderID  string    `bson:"sender_id" json:"sender_id"`
	Content   string    `bson:"content" json:"content"`
	Type      string    `bson:"type" json:"type"` // "text" or "system"
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// ==================== Request Payloads ====================

// CreateRoomPayloadSchema is the JSON body for POST /chat/rooms
type CreateRoomPayloadSchema struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// RoomParamsSchema is the URI param for routes like /chat/rooms/{room_id}
type RoomParamsSchema struct {
	RoomID string `uri:"room_id" binding:"required"`
}

// ==================== WebSocket Messages ====================

// WSMessage is what the CLIENT sends TO the server over WebSocket.
// The "action" field tells the server what the client wants to do.
type WSMessage struct {
	Action  string `json:"action"`            // "send_message", "join_room", "leave_room", "typing"
	RoomID  string `json:"room_id,omitempty"`
	Content string `json:"content,omitempty"`
}

// WSResponse is what the SERVER sends TO the client over WebSocket.
// The "event" field tells the client what happened.
type WSResponse struct {
	Event  string `json:"event"`
	RoomID string      `json:"room_id,omitempty"`
	Data   any         `json:"data,omitempty"`
}
