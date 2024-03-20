package main

import (
	"crypto/md5"
	"encoding/base64"
)

func generateShortCode(input string) (string, error) {

	hash := md5.Sum([]byte(input))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	return encoded[:8], nil
}
