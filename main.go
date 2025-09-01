package main

import (
	"fmt"
	"url/handler"
	"url/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize dependencies
	storageService := service.InitializeStore()
	urlService := service.NewURLService(storageService)

	// Setup middleware for dependency injection
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("urlService", urlService)
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the URL Shortener API",
		})
	})

	r.POST("/create-short-url", func(c *gin.Context) {
		handler.CreateShortUrl(c)
	})

	r.GET("/:shortUrl", func(c *gin.Context) {
		handler.HandleShortUrlRedirect(c)
	})

	err := r.Run(":9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}