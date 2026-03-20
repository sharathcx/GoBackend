package user

import (
	"GoBackend/fastapify"
	"GoBackend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")

	user, err := GetUser(ctx, userID)
	if err != nil {
		statusCode, response := utils.HandleError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, user, "User fetched successfully"))
}

func UpdateUserHandler(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")

	req := fastapify.Req[UpdateUserPayloadSchema](c)

	user, err := UpdateUser(ctx, userID, req)
	if err != nil {
		statusCode, response := utils.HandleError(err)
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, user, "User updated successfully"))
}

func RegisterHandler(c *gin.Context) {
	ctx := c.Request.Context()

	req := fastapify.Req[RegisterPayloadSchema](c)

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		statusCode, response := utils.HandleError(utils.NewApiError(500, err.Error(), utils.ErrInternalError, nil))
		c.JSON(statusCode, response)
		return
	}

	var user User
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
		statusCode, response := utils.HandleError(utils.NewApiError(500, err.Error(), utils.ErrInternalError, nil))
		c.JSON(statusCode, response)
		return
	}

	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, newUser, "User registered successfully"))
}
