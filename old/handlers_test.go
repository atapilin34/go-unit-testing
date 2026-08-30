package old

import (
	"testing"
)

func TestShortHandlerPost(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"валидный HTTP URL", "http://example.com", false},
		{"валидный HTTPS URL", "https://google.com/search?q=test", false},
		{"невалидный URL", "not-a-url", true},
		{"пустая строка", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortener := NewURLShortener()

			shortID, err := shortener.Shorten(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"ошибка = %v, ожидали ошибку = %v",
					err,
					tt.wantErr,
				)

				return
			}

			if tt.wantErr {
				return
			}

			if shortID == "" {
				t.Error("shortID не должен быть пустым")
			}

			if len(shortID) < 6 || len(shortID) > 8 {
				t.Errorf(
					"длина shortID = %d, ожидали от 6 до 8",
					len(shortID),
				)
			}
		})
	}
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
