package main

import (
	"encoding/json"
	"init/example/old"
	"log"
	"net/http"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func main() {
	shortener := old.NewURLShortener()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		var request shortenRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		shortID, err := shortener.Shorten(request.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := shortenResponse{
			ShortURL:    shortID,
			OriginalURL: request.URL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("encode response: %v", err)
		}
	})

	mux.HandleFunc("GET /{shortID}", func(w http.ResponseWriter, r *http.Request) {
		shortID := r.PathValue("shortID")

		originalURL, err := shortener.GetOriginal(shortID)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	})

	log.Println("server started on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
