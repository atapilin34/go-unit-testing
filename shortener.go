package main

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"strings"
	"sync"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLShortener struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

func (us *URLShortener) Shorten(originalURL string) (string, error) {
	originalURL = strings.TrimSpace(originalURL)

	if !isValidURL(originalURL) {
		return "", fmt.Errorf("invalid URL: %s", originalURL)
	}

	us.mu.Lock()
	defer us.mu.Unlock()

	for {
		shortID := generateShortID()
		if shortID == "" {
			return "", fmt.Errorf("failed to generate short ID")
		}

		if _, exists := us.urls[shortID]; exists {
			continue
		}

		us.urls[shortID] = originalURL
		return shortID, nil
	}
}

func (us *URLShortener) GetOriginal(shortID string) (string, error) {
	us.mu.RLock()
	originalURL, exists := us.urls[shortID]
	us.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("original URL not found")
	}

	return originalURL, nil
}

func generateShortID() string {
	const length = 8

	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return ""
	}

	result := make([]byte, length)
	for i, b := range randomBytes {
		result[i] = alphabet[int(b)%len(alphabet)]
	}

	return string(result)
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	return parsedURL.Host != "" && parsedURL.Hostname() != ""
}
