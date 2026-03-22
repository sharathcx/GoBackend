package routes

import (
	"net/http"

	"GoBackend/fastapify"
	"GoBackend/handlers"
	"GoBackend/middleware"
	"GoBackend/schemas"
	"GoBackend/websocket"

	"github.com/gin-gonic/gin"
)

func registerChatRoutes(api *fastapify.Wrapper) {
	go websocket.RunWSHub(websocket.DefaultWSHub)

	chat := api.Group("/chat")

	chat.POST("/rooms", handlers.CreateRoomHandler, middleware.AuthMiddleware()).
		Body(schemas.CreateRoomPayloadSchema{}).
		Response(schemas.RoomSchema{})

	chat.GET("/rooms", handlers.GetUserRoomsHandler, middleware.AuthMiddleware()).
		Response([]schemas.RoomSchema{})

	chat.GET("/rooms/{room_id}/messages", handlers.GetRoomMessagesHandler, middleware.AuthMiddleware()).
		Params(schemas.RoomParamsSchema{}).
		Response([]schemas.MessageSchema{})

	api.Engine.GET("/ws", handlers.ServeWS(websocket.DefaultWSHub))

	api.Engine.GET("/chat", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(websocket.ChatPageHTML))
	})
}
