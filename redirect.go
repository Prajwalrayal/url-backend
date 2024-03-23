package main

func getOrignalURL(shortCode string) (string, error) {
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM urls WHERE short_code = ?", shortCode).Scan(&originalURL)
	return originalURL, err
}
