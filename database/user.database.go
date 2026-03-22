package database

import (
	"GoBackend/schemas"
	"GoBackend/utils"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var userCollection *mongo.Collection

func init() {
	userCollection = OpenCollection("users")
}

func GetUser(ctx context.Context, userID string) (*schemas.UserSchema, *utils.ApiError) {
	var user schemas.UserSchema

	err := userCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, utils.NotFound("user not found")
	}

	return &user, nil
}

func UpdateUser(ctx context.Context, userID string, req *schemas.UpdateUserPayloadSchema) (*schemas.UserSchema, *utils.ApiError) {
	var user schemas.UserSchema
	if req.Email != "" {
		emailExists, err := userCollection.CountDocuments(ctx, bson.M{"email": req.Email})
		if err != nil {
			return nil, utils.InternalError(err.Error())
		}
		if emailExists > 0 {
			return nil, utils.Conflict("email already exists")
		}
	}
	update := bson.M{
		"$set": req,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := userCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userID}, update, opts).Decode(&user)
	if err != nil {
		return nil, utils.NotFound("user not found")
	}

	return &user, nil
}

func InsertUser(ctx context.Context, req *schemas.UserSchema) (*schemas.UserSchema, *utils.ApiError) {
	emailExists, err := userCollection.CountDocuments(ctx, bson.M{"email": req.Email})
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}
	if emailExists > 0 {
		return nil, utils.Conflict("email already exists")
	}
	_, err = userCollection.InsertOne(ctx, req)
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}

	return req, nil
}

func DeleteUser(ctx context.Context, userID string) (*schemas.UserSchema, *utils.ApiError) {
	var user schemas.UserSchema

	err := userCollection.FindOneAndDelete(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, utils.NotFound("user not found")
	}

	return &user, nil
}

func LoginUser(ctx context.Context, req *schemas.UserLoginPayloadSchema) (*schemas.UserSchema, *utils.ApiError) {
	var user schemas.UserSchema

	err := userCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return nil, utils.NotFound("invalid email or password")
	}

	return &user, nil
}

func UpdateAllTokens(ctx context.Context, userID string, accessToken string, refreshToken string) (*schemas.UserSchema, *utils.ApiError) {
	var user schemas.UserSchema
	update := bson.M{
		"$set": bson.M{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"updated_at":    time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := userCollection.FindOneAndUpdate(ctx, bson.M{"user_id": userID}, update, opts).Decode(&user)
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}

	return &user, nil
}
