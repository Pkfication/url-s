package main

import (
	"fmt"
	"url/routes"
	"url/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize dependencies
	storageService := service.InitializeStore()
	urlService := service.NewURLService(storageService)

	// Setup Gin engine
	r := gin.Default()

	// Create services struct for route registration
	services := &routes.Services{
		URLService: urlService,
	}

	// Register all routes with middleware
	routes.RegisterRoutes(r, services)

	// Start the server
	err := r.Run(":9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}