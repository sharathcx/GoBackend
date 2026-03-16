package main

import (
	"fmt"
	"time"

	"GoBackend/fastapify"
	"GoBackend/modules/movie"
	"GoBackend/modules/user"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(fastapify.TimeoutMiddleware(100 * time.Second))

	api := fastapify.New(router)

	movie.RegisterRoutes(api)
	user.RegisterRoutes(api)

	api.SetupSwagger("/openapi.json")

	fmt.Println("Swagger Docs available at: http://localhost:8000/docs")

	if err := router.Run(":8000"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
