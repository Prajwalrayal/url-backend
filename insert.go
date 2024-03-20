package main

func insertURL(originalURL, shortCode string) error {
	_, err := db.Exec("INSERT INTO urls (original_url, short_code) VALUES (?, ?)", originalURL, shortCode)
	return err
}