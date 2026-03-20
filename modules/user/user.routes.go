package user

import (
	"GoBackend/fastapify"
)

func RegisterRoutes(api *fastapify.Wrapper) {
	users := api.Group("/users")

	users.GET("/{user_id}", GetUserHandler).
		Params(UserParamsSchema{}).
		Response(UserSchema{})

	users.PATCH("/{user_id}", UpdateUserHandler).
		Params(UserParamsSchema{}).
		Body(UpdateUserPayloadSchema{}).
		Response(UserSchema{})

	users.POST("", RegisterHandler).
		Body(RegisterPayloadSchema{}).
		Response(UserSchema{})

	users.DELETE("/{user_id}", DeleteUserHandler).
		Params(UserParamsSchema{}).
		Response(UserSchema{})

	users.POST("/login", LoginUserHandler).
		Body(UserLoginPayloadSchema{}).
		Response(UserResponseSchema{})
}
