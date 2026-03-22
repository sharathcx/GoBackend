package schemas

import "github.com/golang-jwt/jwt/v5"

type SignedDetailsSchema struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}
