package user

import (
	"GoBackend/database"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"errors"
)

var userCollection = database.OpenCollection("users")

func GetUser(ctx context.Context, userID string) (*User, error) {
	var user User

	err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(ctx context.Context, userID string, req *UpdateUserPayloadSchema) (*User, error) {
	var user User

	update := bson.M{
		"$set": req,
	}

	//get the updated user
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	
	err := userCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userID}, update, opts).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func InsertUser(ctx context.Context, req *User) (*User, error) {
	emailExists, err := userCollection.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		return nil, err
	}
	if emailExists > 0 {
		return nil, errors.New("email already exists")
	}
	_, err = userCollection.InsertOne(ctx, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}
