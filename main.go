package main

import (
	"fmt"

	"GoBackend/fastapify"
	"GoBackend/modules/movie"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	api := fastapify.New(router)

	movie.RegisterRoutes(api)

	api.SetupSwagger("/openapi.json")

	fmt.Println("Swagger Docs available at: http://localhost:8000/docs")

	if err := router.Run(":8000"); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
