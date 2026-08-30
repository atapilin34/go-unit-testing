package main

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type URLShortener struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		urls: make(map[string]string),
	}
}

// Shorten создает короткий идентификатор для URL
func (us *URLShortener) Shorten(originalURL string) (string, error) {
	// валидация URL
	valid := isValidURL(originalURL)
	if !valid {
		return "", fmt.Errorf("invalid URL %s", originalURL)
	}
	// генерация короткого ID
	us.mu.Lock()
	defer us.mu.Unlock()
	short := generateShortID()
	if short == "" {
		return "", fmt.Errorf("failed to generate short ID")
	}
	// сохранение в map
	us.urls[short] = originalURL
	return short, nil
}

// GetOriginal возвращает оригинальный URL по короткому ID
func (us *URLShortener) GetOriginal(shortID string) (string, error) {
	// TODO: поиск в map
	us.mu.RLock()
	originalURL, ok := us.urls[shortID]
	us.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("original URL not found")
	}
	return originalURL, nil
}

// generateShortID генерирует случайный короткий идентификатор
func generateShortID() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	randomBytes := make([]byte, length)

	if _, err := rand.Read(randomBytes); err != nil {
		return ""
	}

	result := make([]byte, length)

	for i, randomByte := range randomBytes {
		result[i] = alphabet[int(randomByte)%len(alphabet)]
	}

	return string(result)
}

// isValidURL проверяет корректность URL
func isValidURLOld(str string) bool {
	// TODO: валидация URL
	pattern := "^https?://(?:www\\.)?[a-zA-Zа-яА-ЯёЁ0-9-]+\\.[a-zA-Zа-яА-ЯёЁ0-]{2,63}$"
	match, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return match
}

func isValidURL(str string) bool {
	str = strings.TrimSpace(str)

	parsedURL, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	if parsedURL.Host == "" || parsedURL.Hostname() == "" {
		return false
	}

	return true
}
