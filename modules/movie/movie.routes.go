package movie

import (
	"GoBackend/fastapify"
	"GoBackend/middleware/auth"
)

func RegisterRoutes(api *fastapify.Wrapper) {
	movies := api.Group("/movies")

	movies.GET("", GetMoviesHandler).
		Response([]MovieSchema{})

	movies.GET("/{movie_id}", GetMovieHandler, auth.AuthMiddleware()).
		Params(MovieParamsSchema{}).
		Response(MovieSchema{})

	movies.POST("", AddMovieHandler).
		Body(AddMoviePayloadSchema{}).
		Response(MovieSchema{})

	movies.DELETE("/{movie_id}", DeleteMovieHandler).
		Params(MovieParamsSchema{}).
		Response(MovieSchema{})
}
