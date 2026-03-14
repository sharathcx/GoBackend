package movie

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"GoBackend/database"
)

var movieCollection = database.OpenCollection("movies")

func GetMovies(c *gin.Context, req *GetMoviesPayloadSchema) (*[]MovieSchema, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var movies []MovieSchema

	cursor, err := movieCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return &movies, nil
}

func GetMovie(c *gin.Context, req *GetMoviePayloadSchema) (*MovieSchema, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var movie MovieSchema

	err := movieCollection.FindOne(ctx, bson.M{"movie_id": req.MovieID}).Decode(&movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func AddMovie(c *gin.Context, movie *MovieSchema) (*MovieSchema, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	_, err := movieCollection.InsertOne(ctx, movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}


