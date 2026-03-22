package routes

import (
	"GoBackend/fastapify"
	"GoBackend/handlers"
	"GoBackend/schemas"
)

func registerUserRoutes(api *fastapify.Wrapper) {
	users := api.Group("/users")

	users.GET("/{user_id}", handlers.GetUserHandler).
		Params(schemas.UserParamsSchema{}).
		Response(schemas.UserSchema{})

	users.PATCH("/{user_id}", handlers.UpdateUserHandler).
		Params(schemas.UserParamsSchema{}).
		Body(schemas.UpdateUserPayloadSchema{}).
		Response(schemas.UserSchema{})

	users.POST("", handlers.RegisterHandler).
		Body(schemas.RegisterPayloadSchema{}).
		Response(schemas.UserSchema{})

	users.DELETE("/{user_id}", handlers.DeleteUserHandler).
		Params(schemas.UserParamsSchema{}).
		Response(schemas.UserSchema{})

	users.POST("/login", handlers.LoginUserHandler).
		Body(schemas.UserLoginPayloadSchema{}).
		Response(schemas.UserResponseSchema{})
}
