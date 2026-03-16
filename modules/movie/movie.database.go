package movie

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"GoBackend/database"
)

var movieCollection = database.OpenCollection("movies")

func GetMovies(ctx context.Context, req *GetMoviesPayloadSchema) (*[]MovieSchema, error) {
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

func GetMovie(ctx context.Context, movieID string) (*MovieSchema, error) {
	var movie MovieSchema

	err := movieCollection.FindOne(ctx, bson.M{"movie_id": movieID}).Decode(&movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func AddMovie(ctx context.Context, movie *MovieSchema) (*MovieSchema, error) {

	_, err := movieCollection.InsertOne(ctx, movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}
