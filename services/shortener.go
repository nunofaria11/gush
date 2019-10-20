package services

import (
	"gush/models"
	"gush/storage"
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

func createURLInfo(url string, hash string) *models.URLInfo {
	info := models.URLInfo{URL: url, Hash: hash, CreatedAt: time.Now()}
	return &info
}

// GetShortURLInfo Retrieves a URLInfo object
func GetShortURLInfo(hash string) (*models.URLInfo, bool) {

	urlInfo, err := storage.FetchURLInfo(hash)

	if err != nil {
		return nil, false
	}

	return urlInfo, true
}

// SetShortURL Stores a URLInfo object and associates with the hash
func SetShortURL(urlToShorten string) (string, bool) {

	var hash string

	exists := true

	for exists {
		hash = generateHash(6)
		_, exists = GetShortURLInfo(hash)
	}

	urlInfo := createURLInfo(urlToShorten, hash)

	err := storage.StoreURLInfo(urlInfo)
	if err != nil {
		return "", false
	}

	log.Printf("Setting URL info with hash: %v", hash)
	return hash, true
}
