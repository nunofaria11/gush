package shortener

import (
	"gush/entities"
)

// shortURLMap The map used to store URLs in memory
var shortURLMap map[string]*entities.URLInfo

func init() {
	shortURLMap = make(map[string]*entities.URLInfo)
}

// GetShortURL Retrieves a URLInfo object
func GetShortURL(hash string) (*entities.URLInfo, bool) {
	urlInfo, ok := shortURLMap[hash]
	return urlInfo, ok
}

// SetShortURL Stores a URLInfo object and associates with the hash
func SetShortURL(hash string, urlInfo *entities.URLInfo) (*entities.URLInfo, bool) {

	existingURLInfo, exists := shortURLMap[hash]
	if exists {
		return existingURLInfo, false
	}

	shortURLMap[hash] = urlInfo
	return urlInfo, true
}
