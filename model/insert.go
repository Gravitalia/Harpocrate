package model

// URL struct defines how insert of a new short URL should be
type URL struct {
	Id          string  `json:"id"`
	OriginalUrl string  `json:"original_rul"`
	Author      string  `json:"string"`
	Analytics   bool    `json:"analytics"`
	Phishing    float32 `json:"phishing"`
}
