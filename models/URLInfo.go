package models

import "time"

// URLInfo in used to hold URL information
type URLInfo struct {
	URL       string    `json:"url"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}
