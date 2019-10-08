package entities

import "time"

// URLInfo Type used to hold URL info in memory
type URLInfo struct {
	URL  string
	Time time.Time
}

// NewURLInfo Instatiates a new URLInfo object
func NewURLInfo(url string) *URLInfo {
	ui := URLInfo{url, time.Now()}
	return &ui
}
