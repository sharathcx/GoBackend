package chat

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a single WebSocket connection.
// Each user who connects gets one Client struct.
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	UserID   string
	Username string
	Rooms    map[string]bool // which rooms this client has joined (set)
	Send     chan []byte     // outgoing messages queue (buffered)
}

// BroadcastMessage is an instruction to send a message to everyone in a room.
type BroadcastMessage struct {
	RoomID  string
	Message []byte
	Exclude string // skip this userID (the sender doesn't need their own message back)
}

// Hub is the central switchboard. It keeps track of:
// - ALL connected clients (by userID)
// - Which clients are in which rooms
//
// It runs as a single goroutine processing channels — no mutexes needed.
type Hub struct {
	// All connected clients: userID → *Client
	Clients map[string]*Client

	// Room membership: roomID → { userID → *Client }
	Rooms map[string]map[string]*Client

	// Channels — these are the "instructions" sent to the Hub goroutine
	Register   chan *Client          // "a new client connected"
	Unregister chan *Client          // "a client disconnected"
	Broadcast  chan *BroadcastMessage // "send this message to a room"
}

// NewHub creates a Hub with initialized maps and channels.
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[string]map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMessage),
	}
}

// DefaultHub is the singleton — one Hub for the entire application.
var DefaultHub = NewHub()

// Run starts the Hub's event loop. Call this with `go DefaultHub.Run()`.
//
// This is an infinite loop that processes one event at a time:
// - Register: a new client connected → add to Clients map
// - Unregister: a client disconnected → remove from all rooms, close their Send channel
// - Broadcast: someone sent a message → forward to all clients in that room
//
// Because only this one goroutine reads/writes the maps, there are no race conditions.
func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			// If this user was already connected (e.g. page refresh), close the old connection
			if existing, ok := h.Clients[client.UserID]; ok {
				h.removeClient(existing)
			}
			h.Clients[client.UserID] = client
			log.Printf("[Hub] Client registered: %s (%s)", client.Username, client.UserID)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				h.removeClient(client)
				log.Printf("[Hub] Client unregistered: %s (%s)", client.Username, client.UserID)
			}

		case msg := <-h.Broadcast:
			// Get all clients in this room
			room, ok := h.Rooms[msg.RoomID]
			if !ok {
				continue // room doesn't exist in memory, skip
			}

			// Send to every client in the room (except the excluded one)
			for userID, client := range room {
				if userID == msg.Exclude {
					continue
				}
				select {
				case client.Send <- msg.Message:
					// message queued successfully
				default:
					// Send channel is full — this client is too slow.
					// Remove them to prevent blocking the entire broadcast loop.
					h.removeClient(client)
				}
			}
		}
	}
}

// JoinRoom adds a client to a room's in-memory map and notifies others.
func (h *Hub) JoinRoom(client *Client, roomID string) {
	// Create the room map if it doesn't exist yet
	if h.Rooms[roomID] == nil {
		h.Rooms[roomID] = make(map[string]*Client)
	}

	h.Rooms[roomID][client.UserID] = client
	client.Rooms[roomID] = true

	// Notify everyone in the room that this user joined
	response, _ := json.Marshal(WSResponse{
		Event:  "user_joined",
		RoomID: roomID,
		Data: map[string]string{
			"user_id":  client.UserID,
			"username": client.Username,
		},
	})
	h.Broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Message: response,
		Exclude: client.UserID,
	}

	// Send the joining client a list of who's currently online in this room
	h.sendOnlineUsers(client, roomID)
}

// LeaveRoom removes a client from a room and notifies others.
func (h *Hub) LeaveRoom(client *Client, roomID string) {
	room, ok := h.Rooms[roomID]
	if !ok {
		return
	}

	delete(room, client.UserID)
	delete(client.Rooms, roomID)

	// Clean up empty rooms from memory
	if len(room) == 0 {
		delete(h.Rooms, roomID)
	}

	// Notify remaining members
	response, _ := json.Marshal(WSResponse{
		Event:  "user_left",
		RoomID: roomID,
		Data: map[string]string{
			"user_id":  client.UserID,
			"username": client.Username,
		},
	})
	h.Broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Message: response,
		Exclude: client.UserID,
	}
}

// removeClient disconnects a client from everything.
func (h *Hub) removeClient(client *Client) {
	// Leave all rooms this client was in
	for roomID := range client.Rooms {
		h.LeaveRoom(client, roomID)
	}
	close(client.Send)
	client.Conn.Close()
	delete(h.Clients, client.UserID)
}

// sendOnlineUsers sends the list of currently online users in a room to a specific client.
func (h *Hub) sendOnlineUsers(client *Client, roomID string) {
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

	response, _ := json.Marshal(WSResponse{
		Event:  "online_users",
		RoomID: roomID,
		Data:   users,
	})

	// Send directly to this client only, not a broadcast
	client.Send <- response
}
