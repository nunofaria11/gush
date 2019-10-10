package services

import (
	"gush/models"
	"log"
	"math/rand"
	"time"
)

// shortURLMap The map used to store URLs in memory
var shortURLMap map[string]*models.URLInfo
var randomizer *rand.Rand

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
	shortURLMap = make(map[string]*models.URLInfo)
}

func generateHash(strlen int) string {
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[randomizer.Intn(len(chars))]
	}
	return string(result)
}

func createURLInfo(url string) *models.URLInfo {
	ui := models.URLInfo{url, time.Now()}
	return &ui
}

// GetShortURLInfo Retrieves a URLInfo object
func GetShortURLInfo(hash string) (*models.URLInfo, bool) {
	urlInfo, ok := shortURLMap[hash]
	return urlInfo, ok
}

// SetShortURL Stores a URLInfo object and associates with the hash
func SetShortURL(urlToShorten string) (string, bool) {

	var hash string

	urlInfo := createURLInfo(urlToShorten)
	exists := true

	for exists {
		hash = generateHash(8)
		_, exists = shortURLMap[hash]
	}

	shortURLMap[hash] = urlInfo
	log.Printf("Setting URL info with hash: %v", hash)
	return hash, true
}
