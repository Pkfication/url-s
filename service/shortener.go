package service

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// GenerateShortLink creates a short URL from a long URL and user ID
func GenerateShortLink(longURL, userID string) string {
	// Combine URL and user ID for uniqueness
	input := fmt.Sprintf("%s-%s", longURL, userID)
	
	// Generate SHA256 hash
	hash := sha256.Sum256([]byte(input))
	
	// Take first 8 bytes and encode as base64
	shortBytes := hash[:8]
	shortURL := base64.URLEncoding.EncodeToString(shortBytes)
	
	// Remove padding and limit to 8 characters
	shortURL = shortURL[:8]
	
	return shortURL
}
