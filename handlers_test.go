package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShortHandlerPost(t *testing.T) {
	shortener := NewURLShortener()

	req := httptest.NewRequest("POST", "/shorten", nil)

	rec := httptest.NewRecorder()

	if req.Method != http.MethodPost {
		http.NotFound(rec, req)
		return
	}

	shortID, err := shortener.Shorten("https://www.google.com/")
	if err != nil {
		http.Error(rec, err.Error(), http.StatusBadRequest)
		return
	}

	rec.WriteHeader(http.StatusCreated)
	rec.Write([]byte(shortID))
}

func TestShortHandlerGet(t *testing.T) {
	shortener := NewURLShortener()

	originalURL := "https://www.google.com/"

	shortID, err := shortener.Shorten(originalURL)
	if err != nil {
		t.Fatalf("Shorten() ошибка: %v", err)
	}

	gotURL, err := shortener.GetOriginal(shortID)
	if err != nil {
		t.Fatalf("GetOriginal() ошибка: %v", err)
	}

	if gotURL != originalURL {
		t.Errorf(
			"URL = %q, ожидали %q",
			gotURL,
			originalURL,
		)
	}
}
