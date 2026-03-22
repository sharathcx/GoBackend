package schemas

import "time"

type UserParamsSchema struct {
	UserID string `uri:"user_id" binding:"required,min=1"`
}

type UserSchema struct {
	UserID          string        `bson:"user_id" json:"user_id"`
	FirstName       string        `bson:"first_name" json:"first_name"`
	LastName        string        `bson:"last_name" json:"last_name"`
	Email           string        `bson:"email" json:"email"`
	Password        string        `bson:"password" json:"password"`
	Role            string        `bson:"role" json:"role"`
	CreatedAt       time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `bson:"updated_at" json:"updated_at"`
	AccessToken     string        `bson:"access_token" json:"access_token"`
	RefreshToken    string        `bson:"refresh_token" json:"refresh_token"`
	FavouriteGenres []GenreSchema `bson:"favourite_genre" json:"favourite_genre"`
}

type UpdateUserPayloadSchema struct {
	FirstName       string        `bson:"first_name,omitempty" json:"first_name" binding:"omitempty,min=2,max=100"`
	LastName        string        `bson:"last_name,omitempty" json:"last_name" binding:"omitempty,min=2,max=100"`
	Email           string        `bson:"email,omitempty" json:"email" binding:"omitempty,email"`
	Password        string        `bson:"password,omitempty" json:"password" binding:"omitempty,min=2,max=100"`
	Role            string        `bson:"role,omitempty" json:"role" binding:"omitempty,oneof=ADMIN USER"`
	FavouriteGenres []GenreSchema `bson:"favourite_genre,omitempty" json:"favourite_genre" binding:"omitempty,dive"`
	UpdatedAt       time.Time     `bson:"updated_at,omitempty" json:"updated_at" binding:"omitempty"`
}

type RegisterPayloadSchema struct {
	FirstName       string        `json:"first_name" binding:"required,min=2,max=100"`
	LastName        string        `json:"last_name" binding:"required,min=2,max=100"`
	Email           string        `json:"email" binding:"required,email"`
	Password        string        `json:"password" binding:"required,min=2,max=100"`
	Role            string        `json:"role" binding:"required,oneof=ADMIN USER"`
	FavouriteGenres []GenreSchema `json:"favourite_genre" binding:"required,dive"`
}

type UserLoginPayloadSchema struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=2,max=100"`
}

// this is called data transfer object or dto usually contains only the fields that are needed for the response
type UserResponseSchema struct {
	UserID          string        `bson:"user_id" json:"user_id"`
	FirstName       string        `bson:"first_name" json:"first_name"`
	LastName        string        `bson:"last_name" json:"last_name"`
	Email           string        `bson:"email" json:"email"`
	Role            string        `bson:"role" json:"role"`
	FavouriteGenres []GenreSchema `bson:"favourite_genre" json:"favourite_genre"`
	AccessToken     string        `bson:"access_token" json:"access_token"`
	RefreshToken    string        `bson:"refresh_token" json:"refresh_token"`
}
