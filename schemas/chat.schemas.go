package schemas

import (
	"time"

	"github.com/gorilla/websocket"
)

type RoomSchema struct {
	RoomID    string    `bson:"room_id" json:"room_id"`
	Name      string    `bson:"name" json:"name"`
	CreatedBy string    `bson:"created_by" json:"created_by"`
	Members   []string  `bson:"members" json:"members"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type MessageSchema struct {
	MessageID string    `bson:"message_id" json:"message_id"`
	RoomID    string    `bson:"room_id" json:"room_id"`
	SenderID  string    `bson:"sender_id" json:"sender_id"`
	Content   string    `bson:"content" json:"content"`
	Type      string    `bson:"type" json:"type"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type CreateRoomPayloadSchema struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type RoomParamsSchema struct {
	RoomID string `uri:"room_id" binding:"required"`
}

type WSClientSchema struct {
	Hub      *WSHubSchema
	Conn     *websocket.Conn
	UserID   string
	Username string
	Rooms    map[string]bool
	Send     chan []byte
}

type WSBroadcastMessageSchema struct {
	RoomID  string
	Message []byte
	Exclude string
}

type WSHubSchema struct {
	Clients    map[string]*WSClientSchema
	Rooms      map[string]map[string]*WSClientSchema
	Register   chan *WSClientSchema
	Unregister chan *WSClientSchema
	Broadcast  chan *WSBroadcastMessageSchema
}

type WSMessageSchema struct {
	Action  string `json:"action"`
	RoomID  string `json:"room_id,omitempty"`
	Content string `json:"content,omitempty"`
}

type WSResponseSchema struct {
	Event  string `json:"event"`
	RoomID string `json:"room_id,omitempty"`
	Data   any    `json:"data,omitempty"`
}
