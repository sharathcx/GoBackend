package routes

import "GoBackend/fastapify"

func RegisterRoutes(api *fastapify.Wrapper) {
	registerUserRoutes(api)
	registerMovieRoutes(api)
	registerChatRoutes(api)
}
