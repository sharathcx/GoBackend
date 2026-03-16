package database

import (
	"fmt"
	"log"

	"GoBackend/globals"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func init() {

	mongoURI := globals.Vars.MONGO_URI

	if mongoURI == "" {
		log.Fatal("MONGO_URI not set")
	}

	fmt.Println("Mongo URI =", mongoURI)

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}

	Client = client
}

var Client *mongo.Client

func OpenCollection(collectionName string) *mongo.Collection {

	databaseName := globals.Vars.DATABASE_NAME

	if databaseName == "" {
		log.Fatal("DATABASE_NAME not set")
	}

	fmt.Println("Database Name =", databaseName)

	collection := Client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		log.Fatal("Failed to open collection", collectionName)
	}

	return collection
}
