package middleware

import (
	"GoBackend/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetAccessTokenFromHeader(c)
		if err != nil {
			statusCode, response := utils.HandleError(utils.Unauthorized(err.Error()))
			c.JSON(statusCode, response)
			c.Abort()
			return
		}
		if token == "" {
			statusCode, response := utils.HandleError(utils.Unauthorized("token is required"))
			c.JSON(statusCode, response)
			c.Abort()
			return
		}
		claims, err := utils.ValidateToken(token)
		if err != nil {
			statusCode, response := utils.HandleError(utils.Unauthorized(err.Error()))
			c.JSON(statusCode, response)
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("user_id", claims.UserID)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)
		c.Set("role", claims.Role)
		c.Next()
	}
}
