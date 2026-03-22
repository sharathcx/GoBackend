package websocket

import (
	"GoBackend/database"
	"GoBackend/schemas"
	"GoBackend/utils"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 54 * time.Second
	maxMessageSize = 4096
)

func ReadPump(client *schemas.WSClientSchema) {
	defer func() {
		client.Hub.Unregister <- client
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMessage, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg schemas.WSMessageSchema
		if err := json.Unmarshal(rawMessage, &msg); err != nil {
			sendError(client, "invalid message format")
			continue
		}

		handleAction(client, &msg)
	}
}

func WritePump(client *schemas.WSClientSchema) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func handleAction(client *schemas.WSClientSchema, msg *schemas.WSMessageSchema) {
	switch msg.Action {

	case "join_room":
		if msg.RoomID == "" {
			sendError(client, "room_id is required")
			return
		}
		_, err := database.GetRoom(context.Background(), msg.RoomID)
		if err != nil {
			sendError(client, "room not found")
			return
		}
		database.AddMemberToRoom(context.Background(), msg.RoomID, client.UserID)
		JoinRoom(client.Hub, client, msg.RoomID)

	case "leave_room":
		if msg.RoomID == "" {
			sendError(client, "room_id is required")
			return
		}
		database.RemoveMemberFromRoom(context.Background(), msg.RoomID, client.UserID)
		LeaveRoom(client.Hub, client, msg.RoomID)

	case "send_message":
		if msg.RoomID == "" || msg.Content == "" {
			sendError(client, "room_id and content are required")
			return
		}
		if !client.Rooms[msg.RoomID] {
			sendError(client, "you are not a member of this room")
			return
		}

		chatMsg := &schemas.MessageSchema{
			MessageID: utils.InvokeUID("MSG", 4),
			RoomID:    msg.RoomID,
			SenderID:  client.UserID,
			Content:   msg.Content,
			Type:      "text",
			CreatedAt: time.Now(),
		}
		database.InsertMessage(context.Background(), chatMsg)

		response, _ := json.Marshal(schemas.WSResponseSchema{
			Event:  "message",
			RoomID: msg.RoomID,
			Data:   chatMsg,
		})
		client.Hub.Broadcast <- &schemas.WSBroadcastMessageSchema{
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
		response, _ := json.Marshal(schemas.WSResponseSchema{
			Event:  "typing",
			RoomID: msg.RoomID,
			Data: map[string]string{
				"user_id":  client.UserID,
				"username": client.Username,
			},
		})
		client.Hub.Broadcast <- &schemas.WSBroadcastMessageSchema{
			RoomID:  msg.RoomID,
			Message: response,
			Exclude: client.UserID,
		}

	default:
		sendError(client, "unknown action: "+msg.Action)
	}
}

func sendError(client *schemas.WSClientSchema, message string) {
	response, _ := json.Marshal(schemas.WSResponseSchema{
		Event: "error",
		Data:  map[string]string{"message": message},
	})
	select {
	case client.Send <- response:
	default:
		log.Printf("[Chat] Failed to send error to %s: buffer full", client.UserID)
	}
}
