package routes

import (
	"GoBackend/fastapify"
	"GoBackend/handlers"
	"GoBackend/middleware"
	"GoBackend/schemas"
)

func registerMovieRoutes(api *fastapify.Wrapper) {
	movies := api.Group("/movies")

	movies.GET("", handlers.GetMoviesHandler).
		Response([]schemas.MovieSchema{})

	movies.GET("/{movie_id}", handlers.GetMovieHandler, middleware.AuthMiddleware()).
		Params(schemas.MovieParamsSchema{}).
		Response(schemas.MovieSchema{})

	movies.POST("", handlers.AddMovieHandler, middleware.AuthMiddleware()).
		Body(schemas.AddMoviePayloadSchema{}).
		Response(schemas.MovieSchema{})

	movies.DELETE("/{movie_id}", handlers.DeleteMovieHandler, middleware.AuthMiddleware()).
		Params(schemas.MovieParamsSchema{}).
		Response(schemas.MovieSchema{})
}
