package auth

import (
	"GoBackend/globals"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type SignedDetailsSchema struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

var accessTokenSecret = globals.Vars.ACCESS_TOKEN_SECRET
var refreshTokenSecret = globals.Vars.REFRESH_TOKEN_SECRET
var accessTokenExpiryMinutes = globals.Vars.ACCESS_TOKEN_EXPIRY_MINUTES
var refreshTokenExpiryMinutes = globals.Vars.REFRESH_TOKEN_EXPIRY_MINUTES

func GenerateJWT(email, firstName, lastName, role, userID string) (string, string, error) {
	accessTokenClaims := &SignedDetailsSchema{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(accessTokenExpiryMinutes))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "movie-app",
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return "", "", err
	}

	refreshTokenClaims := &SignedDetailsSchema{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(refreshTokenExpiryMinutes))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "movie-app",
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func GetAccessTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}
	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", errors.New("authorization header must start with Bearer")
	}
	tokenString := authHeader[len(prefix):]
	if tokenString == "" {
		return "", errors.New("token is required")
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*SignedDetailsSchema, error) {
	claims := &SignedDetailsSchema{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(accessTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*SignedDetailsSchema); ok && token.Valid {
		return claims, nil
	}
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}
	return nil, errors.New("invalid token")
}
