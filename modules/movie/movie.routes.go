package movie

import (
	"GoBackend/fastapify"
)

func RegisterRoutes(api *fastapify.Wrapper) {
	api.GET("/movies", GetMoviesHandler).
		Response([]MovieSchema{})

	api.GET("/movies/{movie_id}", GetMovieHandler).
		Response(MovieSchema{})

	api.POST("/movies", AddMovieHandler).
		Body(AddMoviePayloadSchema{}).
		Response(MovieSchema{})
}
