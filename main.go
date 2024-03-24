package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

// URL struct to represent a shortened URL
type URL struct {
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
}
type RequestBody struct {
	OriginalURL string `json:"originalURL"`
}

// DB connection variable
var db *sql.DB

func init() {
	// Connect to SQLite database
	var err error
	db, err = sql.Open("sqlite3", "urls.db")
	if err != nil {
		log.Fatal(err)
	}

	// Ensure table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS urls (id INTEGER PRIMARY KEY AUTOINCREMENT, original_url TEXT, short_code TEXT UNIQUE)`)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", shortenHandler)
	mux.HandleFunc("/", redirectHandler)
	handler := cors.AllowAll().Handler(mux)
	fmt.Println("Server started on https://url-backend-2lee.onrender.com/")
	err := http.ListenAndServe(":8081", handler)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	originalURL := requestBody.OriginalURL
	if originalURL == "" {
		http.Error(w, "Missing original URL", http.StatusBadRequest)
		return
	}

	timestamp := time.Now().String()
	shortCode, err := generateShortCode(timestamp)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = insertURL(originalURL, shortCode)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	shortCode = "https://url-backend-2lee.onrender.com/" + shortCode

	w.Header().Set("Content-Type", "application/json")

	response := URL{
		OriginalURL: originalURL,
		ShortCode:   shortCode,
	}
	json.NewEncoder(w).Encode(response)

}
func redirectHandler(w http.ResponseWriter, r *http.Request) {

	shortCode := r.URL.Path[1:]

	originalURL, err := getOrignalURL(shortCode)

	if err != nil {
		return
	}
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
