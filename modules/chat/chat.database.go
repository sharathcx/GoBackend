package chat

import (
	"GoBackend/database"
	"GoBackend/utils"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Two collections — one for rooms, one for messages.
var roomCollection = database.OpenCollection("chat_rooms")
var messageCollection = database.OpenCollection("chat_messages")

// CreateRoom inserts a new room document into MongoDB.
func CreateRoom(ctx context.Context, room *RoomSchema) (*RoomSchema, *utils.ApiError) {
	_, err := roomCollection.InsertOne(ctx, room)
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}
	return room, nil
}

// GetRoom finds a single room by its room_id.
func GetRoom(ctx context.Context, roomID string) (*RoomSchema, *utils.ApiError) {
	var room RoomSchema
	err := roomCollection.FindOne(ctx, bson.M{"room_id": roomID}).Decode(&room)
	if err != nil {
		return nil, utils.NotFound("room not found")
	}
	return &room, nil
}

// GetUserRooms returns all rooms where the given userID is in the members array.
// MongoDB's $in operator isn't needed here — we use a direct match on the array field.
// When you query `bson.M{"members": userID}`, Mongo checks if userID exists anywhere in the array.
func GetUserRooms(ctx context.Context, userID string) ([]RoomSchema, *utils.ApiError) {
	cursor, err := roomCollection.Find(ctx, bson.M{"members": userID})
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}

	var rooms []RoomSchema
	if err := cursor.All(ctx, &rooms); err != nil {
		return nil, utils.InternalError(err.Error())
	}
	return rooms, nil
}

// AddMemberToRoom adds a userID to the room's members array.
// $addToSet ensures no duplicates — if the user is already a member, nothing happens.
func AddMemberToRoom(ctx context.Context, roomID string, userID string) (*RoomSchema, *utils.ApiError) {
	var room RoomSchema
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

// RemoveMemberFromRoom removes a userID from the room's members array.
// $pull removes all occurrences of the value from the array.
func RemoveMemberFromRoom(ctx context.Context, roomID string, userID string) (*RoomSchema, *utils.ApiError) {
	var room RoomSchema
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

// InsertMessage saves a chat message to MongoDB.
func InsertMessage(ctx context.Context, msg *MessageSchema) (*MessageSchema, *utils.ApiError) {
	_, err := messageCollection.InsertOne(ctx, msg)
	if err != nil {
		return nil, utils.InternalError(err.Error())
	}
	return msg, nil
}

// GetMessages returns the most recent messages for a room, newest first.
// limit controls how many messages to return (default 50).
func GetMessages(ctx context.Context, roomID string, limit int64) ([]MessageSchema, *utils.ApiError) {
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

	var messages []MessageSchema
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, utils.InternalError(err.Error())
	}
	return messages, nil
}
