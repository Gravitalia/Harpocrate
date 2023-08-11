package model

// Basics struct represents data to get when
// obtain data from database
type Basics struct {
	OriginalUrl string  `json:"original_url"`
	Phishing    float32 `json:"phishing"`
}
