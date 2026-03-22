package handlers

import (
	"GoBackend/database"
	"GoBackend/fastapify"
	"GoBackend/schemas"
	"GoBackend/utils"
	"GoBackend/websocket"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWS(hub *schemas.WSHubSchema) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			statusCode, response := utils.HandleError(utils.Unauthorized("token query parameter is required"))
			c.JSON(statusCode, response)
			return
		}

		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			statusCode, response := utils.HandleError(utils.Unauthorized("invalid or expired token"))
			c.JSON(statusCode, response)
			return
		}

		conn, wsErr := upgrader.Upgrade(c.Writer, c.Request, nil)
		if wsErr != nil {
			return
		}

		client := &schemas.WSClientSchema{
			Hub:      hub,
			Conn:     conn,
			UserID:   claims.UserID,
			Username: claims.FirstName + " " + claims.LastName,
			Rooms:    make(map[string]bool),
			Send:     make(chan []byte, 256),
		}

		hub.Register <- client

		go websocket.WritePump(client)
		go websocket.ReadPump(client)
	}
}

func CreateRoomHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	req := fastapify.Req[schemas.CreateRoomPayloadSchema](c)

	room := &schemas.RoomSchema{
		RoomID:    utils.InvokeUID("ROM", 4),
		Name:      req.Name,
		CreatedBy: userID,
		Members:   []string{userID},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newRoom, err := database.CreateRoom(ctx, room)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusCreated, newRoom, "Room created successfully")
}

func GetUserRoomsHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	rooms, err := database.GetUserRooms(ctx, userID)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, rooms, "Rooms fetched successfully")
}

func GetRoomMessagesHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")
	roomID := c.Param("room_id")

	room, err := database.GetRoom(ctx, roomID)
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

	messages, err := database.GetMessages(ctx, roomID, 50)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, messages, "Messages fetched successfully")
}
