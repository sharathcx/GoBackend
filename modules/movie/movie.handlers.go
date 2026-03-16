package movie

import (
	"GoBackend/fastapify"
	"GoBackend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMoviesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req GetMoviesPayloadSchema
	if err := c.ShouldBindQuery(&req); err != nil {
		statusCode, response := utils.HandleError(err)
		c.JSON(statusCode, response)
		return
	}

	movies, err := GetMovies(ctx, &req)
	if err != nil {
		statusCode, response := utils.HandleError(utils.NewApiError(500, err.Error(), utils.ErrInternalError, nil))
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, movies, "Movies fetched successfully"))
}

func GetMovieHandler(c *gin.Context) {
	ctx := c.Request.Context()
	movieID := c.Param("movie_id")

	movie, err := GetMovie(ctx, movieID)
	if err != nil {
		statusCode, response := utils.HandleError(utils.NewApiError(500, err.Error(), utils.ErrInternalError, nil))
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, movie, "Movie fetched successfully"))
}

func AddMovieHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req AddMoviePayloadSchema
	if !fastapify.Bind(c, &req) {
		return
	}

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
		statusCode, response := utils.HandleError(utils.NewApiError(500, err.Error(), utils.ErrInternalError, nil))
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, newMovie, "Movie added successfully"))
}
