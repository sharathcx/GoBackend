package main

import (
	"fmt"
	"time"

	"GoBackend/fastapify"
	"GoBackend/globals"
	"GoBackend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(fastapify.TimeoutMiddleware(100 * time.Second))

	api := fastapify.New(router)

	routes.RegisterRoutes(api)

	api.SetupSwagger("/openapi.json")

	fmt.Printf("Swagger Docs available at: http://localhost:%s/docs\n", globals.Vars.PORT)

	if err := router.Run(":" + globals.Vars.PORT); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
