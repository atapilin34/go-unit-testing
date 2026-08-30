package old

import "testing"

func TestURLShortener_Shorten(t *testing.T) {
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

	shortener := NewURLShortener()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortID, err := shortener.Shorten(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ошибка = %v, ожидали ошибку = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(shortID) < 6 {
				t.Errorf("короткий ID слишком короткий: %s", shortID)
			}
		})
	}
}
func TestURLShortener_GetOriginal(t *testing.T) {
	shortener := NewURLShortener()

	originalURL := "https://example.com/very/long/path"

	shortID, err := shortener.Shorten(originalURL)
	if err != nil {
		t.Fatalf("Shorten() вернул ошибку: %v", err)
	}

	gotURL, err := shortener.GetOriginal(shortID)
	if err != nil {
		t.Fatalf("GetOriginal() вернул ошибку: %v", err)
	}

	if gotURL != originalURL {
		t.Errorf(
			"полученный URL = %q, ожидали %q",
			gotURL,
			originalURL,
		)
	}
}
func TestURLShortener_GetOriginalNotFound(t *testing.T) {
	shortener := NewURLShortener()

	_, err := shortener.GetOriginal("unknown1")
	if err == nil {
		t.Error("ожидали ошибку для неизвестного short ID")
	}
}
