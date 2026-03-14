package movie

import (
	"GoBackend/utils"

	"github.com/gin-gonic/gin"
)

func GetMoviesService(c *gin.Context, req *GetMoviesPayloadSchema) (*[]MovieSchema, error) {
	movies, err := GetMovies(c, req)
	if err != nil {
		return nil, utils.NewApiError(500, "failed to fetch movies", utils.ErrInternalError, nil)
	}
	return movies, nil
}

func GetMovieService(c *gin.Context, req *GetMoviePayloadSchema) (*MovieSchema, error) {
	movie, err := GetMovie(c, req)
	if err != nil {
		return nil, utils.NewApiError(500, "failed to fetch movie", utils.ErrInternalError, nil)
	}
	return movie, nil
}

func AddMovieService(c *gin.Context, req *AddMoviePayloadSchema) (*MovieSchema, error) {
	var movie MovieSchema
	movie.MovieID = utils.InvokeUID("MOV", 4)
	movie.Title = req.Title
	movie.PosterPath = req.PosterPath
	movie.YoutubeID = req.YoutubeID
	movie.Genre = req.Genre
	movie.AdminReview = req.AdminReview
	movie.Ranking = req.Ranking

	newMovie, err := AddMovie(c, &movie)
	if err != nil {
		return nil, utils.NewApiError(500, "failed to add movie", utils.ErrInternalError, nil)
	}
	return newMovie, nil
}
