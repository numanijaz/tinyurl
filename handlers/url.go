package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/numanijaz/tinyurl/database"
	"github.com/numanijaz/tinyurl/models"
)

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
	maxAttemptsPerLength := 5
	var _getUniqueHash func(int) string
	_getUniqueHash = func(length int) string {
		if length > 12 {
			log.Printf("Too many hash collisions!!")
			return ""
		}

		var uniqueHash string
		var success bool
		for i := range maxAttemptsPerLength {
			saltedUrl := url
			if i > 0 {
				saltedUrl = fmt.Sprintf("%s-%v", url, i)
			}
			uniqueHash = firstN(getSHA256Hash(saltedUrl), length)

			// DB call to check if new hash already exists
			var existingRecord models.UrlModel
			dbResult := database.DB.First(&existingRecord, "unique_hash = ?", uniqueHash)

			if existingRecord.OriginalUrl == url {
				// hash for an existing URL requested,
				// this may be a security issue because two users requesting
				// short URL for the same url will get the same hash
				success = true
				break
			} else if dbResult.RowsAffected == 0 {
				success = true // unique hash
				break
			}
		}
		if !success {
			log.Printf(
				"Failed to resolve hash collission after %v attempts. Trying with one higher length",
				maxAttemptsPerLength,
			)
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

	shortURL := models.UrlModel{
		OriginalUrl: url,
		UniqueHash:  uniqueHash,
		VisitCount:  0,
		// TODO: set userid that was put in jwt
	}

	created := database.DB.Create(&shortURL)

	print(created)

	c.JSON(http.StatusOK, gin.H{"tinyurl": generateURL(c, uniqueHash)})
}

func GetTinyUrl(c *gin.Context) {
	shortUrlHash := c.Param("tinyurl")

	shortUrl := models.UrlModel{}
	var count int64
	database.DB.First(&shortUrl, "unique_hash = ?", shortUrlHash).Count(&count)
	if count == 0 {
		// The tinyurl handler couldn't find the requested url
		// let the frontend app display error
		c.Redirect(http.StatusSeeOther, "/notfound")
		return
	}

	c.Redirect(http.StatusFound, shortUrl.OriginalUrl)
}
