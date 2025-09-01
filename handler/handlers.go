package handler

import (
	"net/http"
	"url/service"

	"github.com/gin-gonic/gin"
)

// Request model definition
type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId string `json:"user_id" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get URL service from context (injected by middleware)
	urlService := c.MustGet("urlService").(*service.URLService)
	
	shortUrl, err := urlService.CreateShortURL(creationRequest.LongUrl, creationRequest.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	host := "http://localhost:9808/"
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortUrl,
	})
}


func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	
	// Get URL service from context (injected by middleware)
	urlService := c.MustGet("urlService").(*service.URLService)
	
	initialUrl, err := urlService.GetOriginalURL(shortUrl)
	if err != nil || initialUrl == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	c.Redirect(302, initialUrl)
}

// Welcome handles the root endpoint
func Welcome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to the URL Shortener API",
	})
}

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
		"service": "url-shortener",
	})
}

// GetShortUrl retrieves a specific short URL by ID
func GetShortUrl(c *gin.Context) {
	shortUrlID := c.Param("id")
	
	// Get URL service from context (injected by middleware)
	urlService := c.MustGet("urlService").(*service.URLService)
	
	initialUrl, err := urlService.GetOriginalURL(shortUrlID)
	if err != nil || initialUrl == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}
	
	c.JSON(200, gin.H{
		"short_url": shortUrlID,
		"long_url": initialUrl,
	})
}