package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func firstN(s string, n int) string {
	v := []rune(s)
	if n >= len(v) {
		return s
	}
	return string(s[:n])
}

func getSHA256Hash(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	sha := hex.EncodeToString(hasher.Sum(nil))
	return sha
}

func getUniqueHash(url string) string {
	maxAttempts := 5
	var _getUniqueHash func(int) string
	_getUniqueHash = func(length int) string {
		if length > 12 {
			log.Printf("Too many hash collisions!!")
			return ""
		}

		var uniqueHash string
		var success bool
		for i := range maxAttempts {
			saltedUrl := url
			if i > 0 {
				saltedUrl = fmt.Sprintf("%s-%v", url, i)
			}
			uniqueHash = firstN(getSHA256Hash(saltedUrl), length)
			value, exists := db[uniqueHash]
			if exists && value == url {
				// hash for an existing URL requested,
				// this may be a security issue because two users requesting
				// short URL for the same url will get the same hash
				success = true
				break
			} else if !exists {
				success = true // unique hash
				break
			}
		}
		if !success {
			log.Printf("Failed to resolve hash collission after %v attempts.", maxAttempts)
			uniqueHash = _getUniqueHash(length + 1) // try with higher length (recursively)
		}
		return uniqueHash
	}
	return _getUniqueHash(6) // start with length 6
}

func generateURL(c *gin.Context, path string) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	return fmt.Sprintf("%s://%s/%s", scheme, host, strings.TrimPrefix(path, "/"))
}

// The handler that shortens a URL and returns the shortened URL
// which can be navigated later to reach the source URL.
func ShortenUrl(c *gin.Context) {
	url := c.PostForm("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	uniqueHash := getUniqueHash(url)
	if uniqueHash == "" {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Tiny URL could not be created."})
		return
	}

	db[uniqueHash] = url

	c.JSON(http.StatusOK, gin.H{"tinyurl": generateURL(c, uniqueHash)})
}

func GetTinyUrl(c *gin.Context) {
	shortUrl := c.Param("tinyurl")
	if originalUrl, ok := db[shortUrl]; ok {
		c.Redirect(http.StatusFound, originalUrl)
		return
	}
	// The tinyurl handler couldn't find the requested
	// url. Maybe this a known route in frontend app?
	// Serve the frontend app, the app will navigate to notfound
	// page if the route is unknown to the frontend app as well.
	ServeFrontendApp(c)
}
