package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Define the struct wrapper around raw Redis client
type StorageService struct {
	redisClient *redis.Client
}

// Top level declarations for the storeService and Redis context
var (
	storeService = &StorageService{}
    ctx = context.Background()
)

// Note that in a real world usage, the cache duration shouldn't have  
// an expiration time, an LRU policy config should be set where the 
// values that are retrieved less often are purged automatically from 
// the cache and stored back in RDBMS whenever the cache is full

const CacheDuration = 6 * time.Hour

// Initializing the store service and return a store pointer 
func InitializeStore() *StorageService {
	// Get Redis address from environment variable or use default
	redisAddr := "localhost:6379"
	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		redisAddr = addr
	}
	
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error init Redis: %v", err))
	}

	fmt.Printf("\nRedis started successfully: pong message = {%s}", pong)
	storeService.redisClient = redisClient
	return storeService
}

// SaveUrlMapping saves the mapping between short URL and original URL
func (s *StorageService) SaveUrlMapping(shortURL, originalURL, userUUID string) error {
	return s.redisClient.Set(ctx, shortURL, originalURL, CacheDuration).Err()
}

// RetrieveInitialUrl retrieves the original URL from the short URL
func (s *StorageService) RetrieveInitialUrl(shortURL string) (string, error) {
	result, err := s.redisClient.Get(ctx, shortURL).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

// Exists checks if a short URL exists
func (s *StorageService) Exists(shortURL string) bool {
	_, err := s.redisClient.Get(ctx, shortURL).Result()
	return err == nil
}

