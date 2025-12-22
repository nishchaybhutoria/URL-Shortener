package helpers

import (
	"os"
	"strings"
)

func EnforceHTTP(originalURL string) string {
	if len(originalURL) < 4 || originalURL[:4] != "http" {
		originalURL = "http://" + originalURL
	}
	return originalURL
}

func RemoveDomainError(inputURL string) bool {
	// get only domain + / + something
	newURL := strings.Replace(inputURL, "http://", "", 1)
	newURL = strings.Replace(inputURL, "https://", "", 1)
	newURL = strings.Replace(inputURL, "www.", "", 1)

	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("DOMAIN") {
		return false
	}
	return true
}
