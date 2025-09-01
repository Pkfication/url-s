package service

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// URLRepository defines the contract for URL storage operations
type URLRepository interface {
	SaveUrlMapping(shortURL, originalURL, userUUID string) error
	RetrieveInitialUrl(shortURL string) (string, error)
	Exists(shortURL string) bool
}

// RedisRepository implements URLRepository using Redis
type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}
