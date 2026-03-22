package chat

import (
	"GoBackend/fastapify"
	"GoBackend/middleware/auth"
	"GoBackend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader configures the WebSocket upgrade.
// CheckOrigin allows all origins — tighten this in production.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// ServeWS is the WebSocket entry point.
// It's a plain gin.HandlerFunc (NOT a fastapify handler) because:
//   - WebSocket connections are persistent, not request-response
//   - We can't return a value — the connection stays open
//   - It must bypass the timeout middleware
//
// Flow: validate token → upgrade connection → register with Hub → start pumps
func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Get the token from query param
		// WebSocket clients can't set custom headers during the upgrade handshake,
		// so we pass the JWT as a query parameter: ws://host/ws?token=<jwt>
		tokenStr := c.Query("token")
		if tokenStr == "" {
			statusCode, response := utils.HandleError(utils.Unauthorized("token query parameter is required"))
			c.JSON(statusCode, response)
			return
		}

		// Step 2: Validate the JWT (same function the auth middleware uses)
		claims, err := auth.ValidateToken(tokenStr)
		if err != nil {
			statusCode, response := utils.HandleError(utils.Unauthorized("invalid or expired token"))
			c.JSON(statusCode, response)
			return
		}

		// Step 3: Upgrade HTTP → WebSocket
		// After this call, `c.Writer` is no longer an HTTP response writer —
		// it's a WebSocket connection. You can't call c.JSON() after this.
		conn, wsErr := upgrader.Upgrade(c.Writer, c.Request, nil)
		if wsErr != nil {
			return // upgrader already wrote the error response
		}

		// Step 4: Create the client and register with the Hub
		client := &Client{
			Hub:      hub,
			Conn:     conn,
			UserID:   claims.UserID,
			Username: claims.FirstName + " " + claims.LastName,
			Rooms:    make(map[string]bool),
			Send:     make(chan []byte, 256),
		}

		hub.Register <- client

		// Step 5: Start the read and write pumps in separate goroutines.
		// These run for the lifetime of the connection.
		go WritePump(client)
		go ReadPump(client)
	}
}

// ==================== REST Handlers (via fastapify) ====================

// CreateRoomHandler creates a new chat room.
// The authenticated user becomes the first member.
func CreateRoomHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	req := fastapify.Req[CreateRoomPayloadSchema](c)

	room := &RoomSchema{
		RoomID:    utils.InvokeUID("ROM", 4),
		Name:      req.Name,
		CreatedBy: userID,
		Members:   []string{userID},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newRoom, err := CreateRoom(ctx, room)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusCreated, newRoom, "Room created successfully")
}

// GetUserRoomsHandler returns all rooms the authenticated user belongs to.
func GetUserRoomsHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	rooms, err := GetUserRooms(ctx, userID)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, rooms, "Rooms fetched successfully")
}

// GetRoomMessagesHandler returns message history for a room.
// Only members can view messages.
func GetRoomMessagesHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")
	roomID := c.Param("room_id")

	// Verify the user is a member of this room
	room, err := GetRoom(ctx, roomID)
	if err != nil {
		return err
	}

	isMember := false
	for _, member := range room.Members {
		if member == userID {
			isMember = true
			break
		}
	}
	if !isMember {
		return utils.Forbidden("you are not a member of this room")
	}

	messages, err := GetMessages(ctx, roomID, 50)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, messages, "Messages fetched successfully")
}
