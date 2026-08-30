package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortHandlerPost(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "валидный HTTP URL",
			body:       `{"url":"http://example.com"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "валидный HTTPS URL",
			body:       `{"url":"https://google.com/search?q=test"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "невалидный URL",
			body:       `{"url":"not-a-url"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "пустая строка",
			body:       `{"url":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "некорректный JSON",
			body:       `{"url":`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := NewURLShortener()
			handler := setupHandlers(shortener)

			req := httptest.NewRequest(
				http.MethodPost,
				"/shorten",
				bytes.NewBufferString(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("статус = %d, ожидали %d", rec.Code, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusCreated {
				return
			}

			var response shortenResponse
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("ошибка чтения ответа: %v", err)
			}

			if response.ShortURL == "" {
				t.Error("short_url не должен быть пустым")
			}
		})
	}
}

func TestShortHandlerGet(t *testing.T) {
	shortener := NewURLShortener()
	handler := setupHandlers(shortener)
	originalURL := "https://www.google.com/"

	shortID, err := shortener.Shorten(originalURL)
	if err != nil {
		t.Fatalf("Shorten() ошибка: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/"+shortID, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("статус = %d, ожидали %d", rec.Code, http.StatusFound)
	}

	if location := rec.Header().Get("Location"); location != originalURL {
		t.Errorf("Location = %q, ожидали %q", location, originalURL)
	}
}

func TestShortHandlerGetNotFound(t *testing.T) {
	shortener := NewURLShortener()
	handler := setupHandlers(shortener)

	req := httptest.NewRequest(http.MethodGet, "/unknown1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("статус = %d, ожидали %d", rec.Code, http.StatusNotFound)
	}
}
