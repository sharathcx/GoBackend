package movie

import (
	"GoBackend/fastapify"
)

func RegisterRoutes(api *fastapify.Wrapper) {
	fastapify.Get(api, "/movies", GetMoviesService)
	fastapify.Get(api, "/movies/{movie_id}", GetMovieService)
	fastapify.Post(api, "/movies", AddMovieService)
}
