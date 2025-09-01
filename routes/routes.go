package routes

import (
	"url/handler"
	"url/middleware"
	"url/service"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all application routes with middleware
func RegisterRoutes(r *gin.Engine, services *Services) {
	// Create route registry with default configurations
	routeRegistry := middleware.DefaultRouteRegistry()
	
	// Apply route-specific middleware
	r.Use(routeRegistry.RouteMiddleware())
	
	// Apply dependency injection middleware
	r.Use(func(c *gin.Context) {
		c.Set("urlService", services.URLService)
		c.Next()
	})
	
	// Register all route groups
	registerHealthRoutes(r)
	registerAPIRoutes(r)
	registerRedirectRoutes(r)
}

// Services holds all service dependencies
type Services struct {
	URLService *service.URLService
}

// registerHealthRoutes registers health and status endpoints
func registerHealthRoutes(r *gin.Engine) {
	// Welcome endpoint
	r.GET("/", handler.Welcome)
	
	// Health check endpoint
	r.GET("/health", handler.HealthCheck)
}

// registerAPIRoutes registers API endpoints
func registerAPIRoutes(r *gin.Engine) {
	// API v1 routes
	api := r.Group("/api/v1")
	{
		// URL shortening endpoints
		api.POST("/short-urls", handler.CreateShortUrl)
		api.GET("/short-urls/:id", handler.GetShortUrl)
		
		// Analytics endpoints (future)
		// api.GET("/analytics/:shortUrl", handler.GetAnalytics)
		// api.GET("/user/:userId/urls", handler.GetUserUrls)
	}
}

// registerRedirectRoutes registers redirect endpoints
func registerRedirectRoutes(r *gin.Engine) {
	// Short URL redirect endpoint
	r.GET("/:shortUrl", handler.HandleShortUrlRedirect)
}
