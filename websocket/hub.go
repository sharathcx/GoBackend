package websocket

import (
	"GoBackend/schemas"
	"encoding/json"
	"log"
)

func NewWSHub() *schemas.WSHubSchema {
	return &schemas.WSHubSchema{
		Clients:    make(map[string]*schemas.WSClientSchema),
		Rooms:      make(map[string]map[string]*schemas.WSClientSchema),
		Register:   make(chan *schemas.WSClientSchema),
		Unregister: make(chan *schemas.WSClientSchema),
		Broadcast:  make(chan *schemas.WSBroadcastMessageSchema),
	}
}

var DefaultWSHub = NewWSHub()

func RunWSHub(h *schemas.WSHubSchema) {
	for {
		select {

		case client := <-h.Register:
			if existing, ok := h.Clients[client.UserID]; ok {
				removeClient(h, existing)
			}
			h.Clients[client.UserID] = client
			log.Printf("[Hub] Client registered: %s (%s)", client.Username, client.UserID)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				removeClient(h, client)
				log.Printf("[Hub] Client unregistered: %s (%s)", client.Username, client.UserID)
			}

		case msg := <-h.Broadcast:
			room, ok := h.Rooms[msg.RoomID]
			if !ok {
				continue
			}

			for userID, client := range room {
				if userID == msg.Exclude {
					continue
				}
				select {
				case client.Send <- msg.Message:
				default:
					removeClient(h, client)
				}
			}
		}
	}
}

func JoinRoom(h *schemas.WSHubSchema, client *schemas.WSClientSchema, roomID string) {
	if h.Rooms[roomID] == nil {
		h.Rooms[roomID] = make(map[string]*schemas.WSClientSchema)
	}

	h.Rooms[roomID][client.UserID] = client
	client.Rooms[roomID] = true

	response, _ := json.Marshal(schemas.WSResponseSchema{
		Event:  "user_joined",
		RoomID: roomID,
		Data: map[string]string{
			"user_id":  client.UserID,
			"username": client.Username,
		},
	})
	h.Broadcast <- &schemas.WSBroadcastMessageSchema{
		RoomID:  roomID,
		Message: response,
		Exclude: client.UserID,
	}

	sendOnlineUsers(h, client, roomID)
}

func LeaveRoom(h *schemas.WSHubSchema, client *schemas.WSClientSchema, roomID string) {
	room, ok := h.Rooms[roomID]
	if !ok {
		return
	}

	delete(room, client.UserID)
	delete(client.Rooms, roomID)

	if len(room) == 0 {
		delete(h.Rooms, roomID)
	}

	response, _ := json.Marshal(schemas.WSResponseSchema{
		Event:  "user_left",
		RoomID: roomID,
		Data: map[string]string{
			"user_id":  client.UserID,
			"username": client.Username,
		},
	})
	h.Broadcast <- &schemas.WSBroadcastMessageSchema{
		RoomID:  roomID,
		Message: response,
		Exclude: client.UserID,
	}
}

func removeClient(h *schemas.WSHubSchema, client *schemas.WSClientSchema) {
	for roomID := range client.Rooms {
		LeaveRoom(h, client, roomID)
	}
	close(client.Send)
	client.Conn.Close()
	delete(h.Clients, client.UserID)
}

func sendOnlineUsers(h *schemas.WSHubSchema, client *schemas.WSClientSchema, roomID string) {
	room, ok := h.Rooms[roomID]
	if !ok {
		return
	}

	users := []map[string]string{}
	for _, c := range room {
		users = append(users, map[string]string{
			"user_id":  c.UserID,
			"username": c.Username,
		})
	}

	response, _ := json.Marshal(schemas.WSResponseSchema{
		Event:  "online_users",
		RoomID: roomID,
		Data:   users,
	})

	client.Send <- response
}
