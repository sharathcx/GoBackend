package database

import (
	"GoBackend/schemas"
	"GoBackend/utils"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection

func init() {
	movieCollection = OpenCollection("movies")
}

func GetMovies(ctx context.Context, req *schemas.GetMoviesPayloadSchema) (*[]schemas.MovieSchema, *utils.ApiError) {
	var movies []schemas.MovieSchema

	cursor, err := movieCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &movies); err != nil {
		return nil, utils.InternalError(err.Error())
	}

	return &movies, nil
}

func GetMovie(ctx context.Context, movieID string) (*schemas.MovieSchema, *utils.ApiError) {
	var movie schemas.MovieSchema

	err := movieCollection.FindOne(ctx, bson.M{"movie_id": movieID}).Decode(&movie)
	if err != nil {
		return nil, utils.NotFound("movie not found")
	}

	return &movie, nil
}

func AddMovie(ctx context.Context, movie *schemas.MovieSchema) (*schemas.MovieSchema, *utils.ApiError) {
	_, err := movieCollection.InsertOne(ctx, movie)
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}

	return movie, nil
}

func DeleteMovie(ctx context.Context, movieID string) (*schemas.MovieSchema, *utils.ApiError) {
	var movie schemas.MovieSchema

	err := movieCollection.FindOneAndDelete(ctx, bson.M{"movie_id": movieID}).Decode(&movie)
	if err != nil {
		return nil, utils.NotFound("movie not found")
	}

	return &movie, nil
}
