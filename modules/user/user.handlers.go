package user

import (
	"GoBackend/fastapify"
	"GoBackend/middleware/auth"
	"GoBackend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	params := fastapify.Params[UserParamsSchema](c)

	user, err := GetUser(ctx, params.UserID)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, user, "User fetched successfully")
}

func UpdateUserHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	params := fastapify.Params[UserParamsSchema](c)

	req := fastapify.Req[UpdateUserPayloadSchema](c)
	req.UpdatedAt = time.Now()
	user, err := UpdateUser(ctx, params.UserID, req)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, user, "User updated successfully")
}

func RegisterHandler(c *gin.Context) any {
	ctx := c.Request.Context()

	req := fastapify.Req[RegisterPayloadSchema](c)

	hashedPassword, hashErr := HashPassword(req.Password)
	if hashErr != nil {
		return utils.InternalError(hashErr.Error())
	}

	var user UserSchema
	user.UserID = utils.InvokeUID("USR", 4)
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.Password = hashedPassword
	user.Role = req.Role
	user.FavouriteGenres = req.FavouriteGenres
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	newUser, err := InsertUser(ctx, &user)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, newUser, "User registered successfully")
}

func DeleteUserHandler(c *gin.Context) any {
	ctx := c.Request.Context()
	params := fastapify.Params[UserParamsSchema](c)

	user, err := DeleteUser(ctx, params.UserID)
	if err != nil {
		return err
	}

	return utils.NewApiResponse(http.StatusOK, user, "User deleted successfully")
}

func LoginUserHandler(c *gin.Context) any {
	ctx := c.Request.Context()

	req := fastapify.Req[UserLoginPayloadSchema](c)

	foundUser, err := LoginUser(ctx, req)
	if err != nil {
		return err
	}
	verifyErr := VerifyPassword(req.Password, foundUser.Password)
	if verifyErr != nil {
		return utils.Unauthorized("invalid email or password")
	}

	accessToken, refreshToken, jwtErr := auth.GenerateJWT(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)
	if jwtErr != nil {
		return utils.InternalError(jwtErr.Error())
	}

	_, tokenErr := UpdateAllTokens(ctx, foundUser.UserID, accessToken, refreshToken)
	if tokenErr != nil {
		return tokenErr
	}

	response := UserResponseSchema{
		UserID:          foundUser.UserID,
		FirstName:       foundUser.FirstName,
		LastName:        foundUser.LastName,
		Email:           foundUser.Email,
		Role:            foundUser.Role,
		FavouriteGenres: foundUser.FavouriteGenres,
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
	}

	return utils.NewApiResponse(http.StatusOK, response, "User logged in successfully")
}
