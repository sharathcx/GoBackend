package user

import (
	"GoBackend/fastapify"
)

func RegisterRoutes(api *fastapify.Wrapper) {
	api.GET("/users/{user_id}", GetUserHandler).
		Response(User{})

	api.PATCH("/users/{user_id}", UpdateUserHandler).
		Body(UpdateUserPayloadSchema{}).
		Response(User{})

	api.POST("/users", RegisterHandler).
		Body(RegisterPayloadSchema{}).
		Response(User{})
}
