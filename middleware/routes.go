package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RouteConfig holds configuration for specific routes
type RouteConfig struct {
	SkipAuth     bool
	RateLimit    bool
	Logging      bool
	Metrics      bool
	CacheControl string
}

// RouteRegistry manages route configurations
type RouteRegistry struct {
	routes map[string]*RouteConfig
}

// NewRouteRegistry creates a new route registry
func NewRouteRegistry() *RouteRegistry {
	return &RouteRegistry{
		routes: make(map[string]*RouteConfig),
	}
}

// RegisterRoute configures a route with specific middleware settings
func (r *RouteRegistry) RegisterRoute(method, path string, config *RouteConfig) {
	key := method + ":" + path
	r.routes[key] = config
}

// GetRouteConfig retrieves configuration for a specific route
func (r *RouteRegistry) GetRouteConfig(method, path string) *RouteConfig {
	// Try exact match first
	key := method + ":" + path
	if config, exists := r.routes[key]; exists {
		return config
	}
	
	// Try pattern matching for dynamic routes
	for routeKey, config := range r.routes {
		if r.matchesPattern(routeKey, method+":"+path) {
			return config
		}
	}
	
	// Return default config
	return &RouteConfig{
		SkipAuth:     false,
		RateLimit:    true,
		Logging:      true,
		Metrics:      true,
		CacheControl: "no-cache",
	}
}

// matchesPattern checks if a request path matches a registered route pattern
func (r *RouteRegistry) matchesPattern(pattern, request string) bool {
	patternParts := strings.Split(pattern, ":")
	requestParts := strings.Split(request, ":")
	
	if len(patternParts) != 2 || len(requestParts) != 2 {
		return false
	}
	
	// Check method
	if patternParts[0] != requestParts[0] {
		return false
	}
	
	// Check path pattern
	return r.pathMatches(patternParts[1], requestParts[1])
}

// pathMatches checks if a request path matches a route pattern
func (r *RouteRegistry) pathMatches(pattern, request string) bool {
	patternSegments := strings.Split(pattern, "/")
	requestSegments := strings.Split(request, "/")
	
	if len(patternSegments) != len(requestSegments) {
		return false
	}
	
	for i, patternSeg := range patternSegments {
		if strings.HasPrefix(patternSeg, ":") {
			// Dynamic segment, always matches
			continue
		}
		if patternSeg != requestSegments[i] {
			return false
		}
	}
	
	return true
}

// RouteMiddleware applies route-specific middleware based on configuration
func (r *RouteRegistry) RouteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := r.GetRouteConfig(c.Request.Method, c.Request.URL.Path)
		
		// Set route configuration in context
		c.Set("routeConfig", config)
		
		// Apply route-specific middleware
		if config.SkipAuth {
			c.Set("skipAuth", true)
		}
		
		if config.RateLimit {
			c.Set("rateLimit", true)
		}
		
		if config.Logging {
			c.Set("enableLogging", true)
		}
		
		if config.Metrics {
			c.Set("enableMetrics", true)
		}
		
		if config.CacheControl != "" {
			c.Header("Cache-Control", config.CacheControl)
		}
		
		c.Next()
	}
}

// GetRateLimitConfig returns rate limiting configuration for a route
func (r *RouteRegistry) GetRateLimitConfig(method, path string) *RateLimitConfig {
	config := r.GetRouteConfig(method, path)
	if !config.RateLimit {
		return nil
	}
	
	// Return default rate limit config
	return &RateLimitConfig{
		Window: time.Minute,  // 1 minute window
		Limit:  100,          // 100 requests per minute
	}
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Window time.Duration
	Limit  int
}

// DefaultRouteRegistry creates a registry with common route configurations
func DefaultRouteRegistry() *RouteRegistry {
	registry := NewRouteRegistry()
	
	// Health and status routes (no auth, no rate limiting)
	registry.RegisterRoute("GET", "/", &RouteConfig{
		SkipAuth:     true,
		RateLimit:    false,
		Logging:      true,
		Metrics:      true,
		CacheControl: "public, max-age=300",
	})
	
	registry.RegisterRoute("GET", "/health", &RouteConfig{
		SkipAuth:     true,
		RateLimit:    false,
		Logging:      false,
		Metrics:      true,
		CacheControl: "no-cache",
	})
	
	// API routes (with auth and rate limiting)
	registry.RegisterRoute("POST", "/api/v1/short-urls", &RouteConfig{
		SkipAuth:     false,
		RateLimit:    true,
		Logging:      true,
		Metrics:      true,
		CacheControl: "no-cache",
	})
	
	// GetShortUrl with stricter rate limiting (more sensitive endpoint)
	registry.RegisterRoute("GET", "/api/v1/short-urls/:id", &RouteConfig{
		SkipAuth:     false,
		RateLimit:    true,
		Logging:      true,
		Metrics:      true,
		CacheControl: "public, max-age=3600",
	})
	
	// Redirect routes (public, with rate limiting)
	registry.RegisterRoute("GET", "/:shortUrl", &RouteConfig{
		SkipAuth:     true,
		RateLimit:    true,
		Logging:      true,
		Metrics:      true,
		CacheControl: "no-cache",
	})
	
	return registry
}
