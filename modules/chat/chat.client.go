package chat

import (
	"GoBackend/utils"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Timing constants for WebSocket keepalive.
//
// WebSocket connections can silently die (network change, laptop sleep, etc).
// Ping/pong is how we detect dead connections:
//   - Server sends a ping every 54s (pingPeriod)
//   - Client must respond with a pong within 60s (pongWait)
//   - If no pong arrives, the read deadline expires and ReadPump exits
const (
	writeWait      = 10 * time.Second  // max time to write a message
	pongWait       = 60 * time.Second  // max time to wait for pong from client
	pingPeriod     = 54 * time.Second  // how often to send pings (must be < pongWait)
	maxMessageSize = 4096              // max message size in bytes
)

// ReadPump runs in its own goroutine. It reads messages from the WebSocket
// connection and dispatches them based on the "action" field.
//
// When this function returns (connection closed, error, timeout),
// the client is unregistered from the Hub.
func ReadPump(client *Client) {
	// When ReadPump exits (for any reason), clean up
	defer func() {
		client.Hub.Unregister <- client
	}()

	// Configure the connection
	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))

	// When we receive a pong, reset the read deadline.
	// This keeps the connection alive as long as the client is responsive.
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Main read loop — runs until the connection breaks
	for {
		_, rawMessage, err := client.Conn.ReadMessage()
		if err != nil {
			// Connection closed or errored — exit the loop
			break
		}

		// Parse the JSON message from the client
		var msg WSMessage
		if err := json.Unmarshal(rawMessage, &msg); err != nil {
			sendError(client, "invalid message format")
			continue
		}

		// Dispatch based on what the client wants to do
		handleAction(client, &msg)
	}
}

// WritePump runs in its own goroutine. It sends messages from the
// Send channel to the WebSocket connection, and sends periodic pings.
func WritePump(client *Client) {
	// Ticker sends a ping every 54 seconds to keep the connection alive
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			// Set a deadline — if we can't write within 10s, the connection is stuck
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// Hub closed the Send channel → connection is being terminated
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write the message to the WebSocket
			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			// Send a ping to check if the client is still alive
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleAction dispatches a client's message to the right handler.
func handleAction(client *Client, msg *WSMessage) {
	switch msg.Action {

	case "join_room":
		if msg.RoomID == "" {
			sendError(client, "room_id is required")
			return
		}
		// Verify the room exists in the database before joining
		_, err := GetRoom(context.Background(), msg.RoomID)
		if err != nil {
			sendError(client, "room not found")
			return
		}
		// Add to DB members list
		AddMemberToRoom(context.Background(), msg.RoomID, client.UserID)
		// Add to Hub's in-memory room map
		client.Hub.JoinRoom(client, msg.RoomID)

	case "leave_room":
		if msg.RoomID == "" {
			sendError(client, "room_id is required")
			return
		}
		RemoveMemberFromRoom(context.Background(), msg.RoomID, client.UserID)
		client.Hub.LeaveRoom(client, msg.RoomID)

	case "send_message":
		if msg.RoomID == "" || msg.Content == "" {
			sendError(client, "room_id and content are required")
			return
		}
		// Check the client is actually in this room
		if !client.Rooms[msg.RoomID] {
			sendError(client, "you are not a member of this room")
			return
		}

		// Create and persist the message
		chatMsg := &MessageSchema{
			MessageID: utils.InvokeUID("MSG", 4),
			RoomID:    msg.RoomID,
			SenderID:  client.UserID,
			Content:   msg.Content,
			Type:      "text",
			CreatedAt: time.Now(),
		}
		InsertMessage(context.Background(), chatMsg)

		// Broadcast to everyone in the room (including sender)
		response, _ := json.Marshal(WSResponse{
			Event:  "message",
			RoomID: msg.RoomID,
			Data:   chatMsg,
		})
		client.Hub.Broadcast <- &BroadcastMessage{
			RoomID:  msg.RoomID,
			Message: response,
		}

	case "typing":
		if msg.RoomID == "" {
			return
		}
		if !client.Rooms[msg.RoomID] {
			return
		}
		response, _ := json.Marshal(WSResponse{
			Event:  "typing",
			RoomID: msg.RoomID,
			Data: map[string]string{
				"user_id":  client.UserID,
				"username": client.Username,
			},
		})
		client.Hub.Broadcast <- &BroadcastMessage{
			RoomID:  msg.RoomID,
			Message: response,
			Exclude: client.UserID, // don't tell the typer they're typing
		}

	default:
		sendError(client, "unknown action: "+msg.Action)
	}
}

// sendError sends an error message directly to a single client.
func sendError(client *Client, message string) {
	response, _ := json.Marshal(WSResponse{
		Event: "error",
		Data:  map[string]string{"message": message},
	})
	select {
	case client.Send <- response:
	default:
		log.Printf("[Chat] Failed to send error to %s: buffer full", client.UserID)
	}
}
