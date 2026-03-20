package movie

import (
	"GoBackend/fastapify"
	"GoBackend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMoviesHandler(c *gin.Context) any {
	ctx := c.Request.Context()

	req := fastapify.Req[GetMoviesPayloadSchema](c)

	movies, err := GetMovies(ctx, req)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, movies, "Movies fetched successfully")
}

func GetMovieHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	params := fastapify.Params[MovieParamsSchema](c)

	movie, err := GetMovie(ctx, params.MovieID)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, movie, "Movie fetched successfully")
}

func AddMovieHandler(c *gin.Context) any {
	ctx := c.Request.Context()

	req := fastapify.Req[AddMoviePayloadSchema](c)

	var movie MovieSchema
	movie.MovieID = utils.InvokeUID("MOV", 4)
	movie.Title = req.Title
	movie.PosterPath = req.PosterPath
	movie.YoutubeID = req.YoutubeID
	movie.Genre = req.Genre
	movie.AdminReview = req.AdminReview
	movie.Ranking = req.Ranking

	newMovie, err := AddMovie(ctx, &movie)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, newMovie, "Movie added successfully")
}

func DeleteMovieHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	params := fastapify.Params[MovieParamsSchema](c)

	movie, err := DeleteMovie(ctx, params.MovieID)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, movie, "Movie deleted successfully")
}
