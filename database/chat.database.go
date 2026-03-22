package database

import (
	"GoBackend/schemas"
	"GoBackend/utils"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var roomCollection *mongo.Collection
var messageCollection *mongo.Collection

func init() {
	roomCollection = OpenCollection("chat_rooms")
	messageCollection = OpenCollection("chat_messages")
}

func CreateRoom(ctx context.Context, room *schemas.RoomSchema) (*schemas.RoomSchema, *utils.ApiError) {
	_, err := roomCollection.InsertOne(ctx, room)
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}
	return room, nil
}

func GetRoom(ctx context.Context, roomID string) (*schemas.RoomSchema, *utils.ApiError) {
	var room schemas.RoomSchema
	err := roomCollection.FindOne(ctx, bson.M{"room_id": roomID}).Decode(&room)
	if err != nil {
		return nil, utils.NotFound("room not found")
	}
	return &room, nil
}

func GetUserRooms(ctx context.Context, userID string) ([]schemas.RoomSchema, *utils.ApiError) {
	cursor, err := roomCollection.Find(ctx, bson.M{"members": userID})
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}

	var rooms []schemas.RoomSchema
	if err := cursor.All(ctx, &rooms); err != nil {
		return nil, utils.InternalError(err.Error())
	}
	return rooms, nil
}

func AddMemberToRoom(ctx context.Context, roomID string, userID string) (*schemas.RoomSchema, *utils.ApiError) {
	var room schemas.RoomSchema
	update := bson.M{
		"$addToSet": bson.M{"members": userID},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := roomCollection.FindOneAndUpdate(ctx, bson.M{"room_id": roomID}, update, opts).Decode(&room)
	if err != nil {
		return nil, utils.NotFound("room not found")
	}
	return &room, nil
}

func RemoveMemberFromRoom(ctx context.Context, roomID string, userID string) (*schemas.RoomSchema, *utils.ApiError) {
	var room schemas.RoomSchema
	update := bson.M{
		"$pull": bson.M{"members": userID},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := roomCollection.FindOneAndUpdate(ctx, bson.M{"room_id": roomID}, update, opts).Decode(&room)
	if err != nil {
		return nil, utils.NotFound("room not found")
	}
	return &room, nil
}

func InsertMessage(ctx context.Context, msg *schemas.MessageSchema) (*schemas.MessageSchema, *utils.ApiError) {
	_, err := messageCollection.InsertOne(ctx, msg)
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}
	return msg, nil
}

func GetMessages(ctx context.Context, roomID string, limit int64) ([]schemas.MessageSchema, *utils.ApiError) {
	if limit <= 0 {
		limit = 50
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(limit)

	cursor, err := messageCollection.Find(ctx, bson.M{"room_id": roomID}, opts)
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}

	var messages []schemas.MessageSchema
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, utils.InternalError(err.Error())
	}
	return messages, nil
}
