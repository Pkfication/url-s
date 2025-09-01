package service

// URLService handles URL shortening business logic
type URLService struct {
	urlRepo URLRepository
}

// NewURLService creates a new URL service with dependency injection
func NewURLService(urlRepo URLRepository) *URLService {
	return &URLService{
		urlRepo: urlRepo,
	}
}

// CreateShortURL creates a short URL and stores the mapping
func (s *URLService) CreateShortURL(longURL, userID string) (string, error) {
	shortURL := GenerateShortLink(longURL, userID)
	
	err := s.urlRepo.SaveUrlMapping(shortURL, longURL, userID)
	if err != nil {
		return "", err
	}
	
	return shortURL, nil
}

// GetOriginalURL retrieves the original URL from a short URL
func (s *URLService) GetOriginalURL(shortURL string) (string, error) {
	return s.urlRepo.RetrieveInitialUrl(shortURL)
}
