package chat

import (
	"net/http"

	"GoBackend/fastapify"
	"GoBackend/middleware/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *fastapify.Wrapper) {
	// Start the Hub goroutine — this runs for the lifetime of the server
	go DefaultHub.Run()

	// REST endpoints via fastapify (room management + message history)
	chat := api.Group("/chat")

	chat.POST("/rooms", CreateRoomHandler, auth.AuthMiddleware()).
		Body(CreateRoomPayloadSchema{}).
		Response(RoomSchema{})

	chat.GET("/rooms", GetUserRoomsHandler, auth.AuthMiddleware()).
		Response([]RoomSchema{})

	chat.GET("/rooms/{room_id}/messages", GetRoomMessagesHandler, auth.AuthMiddleware()).
		Params(RoomParamsSchema{}).
		Response([]MessageSchema{})

	// WebSocket endpoint — registered directly on the Gin engine.
	// This bypasses fastapify (no auto-serialization) and the timeout middleware.
	api.Engine.GET("/ws", ServeWS(DefaultHub))

	// Chat test UI
	api.Engine.GET("/chat", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(chatPageHTML))
	})
}
