package handlers

import (
	"GoBackend/database"
	"GoBackend/fastapify"
	"GoBackend/schemas"
	"GoBackend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMoviesHandler(c *gin.Context) any {
	ctx := c.Request.Context()

	req := fastapify.Req[schemas.GetMoviesPayloadSchema](c)

	movies, err := database.GetMovies(ctx, req)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, movies, "Movies fetched successfully")
}

func GetMovieHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	params := fastapify.Params[schemas.MovieParamsSchema](c)

	movie, err := database.GetMovie(ctx, params.MovieID)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, movie, "Movie fetched successfully")
}

func AddMovieHandler(c *gin.Context) any {
	ctx := c.Request.Context()

	req := fastapify.Req[schemas.AddMoviePayloadSchema](c)

	var movie schemas.MovieSchema
	movie.MovieID = utils.InvokeUID("MOV", 4)
	movie.Title = req.Title
	movie.PosterPath = req.PosterPath
	movie.YoutubeID = req.YoutubeID
	movie.Genre = req.Genre
	movie.AdminReview = req.AdminReview
	movie.Ranking = req.Ranking

	newMovie, err := database.AddMovie(ctx, &movie)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, newMovie, "Movie added successfully")
}

func DeleteMovieHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	params := fastapify.Params[schemas.MovieParamsSchema](c)

	movie, err := database.DeleteMovie(ctx, params.MovieID)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, movie, "Movie deleted successfully")
}
