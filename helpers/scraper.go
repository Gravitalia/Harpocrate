package helpers

import (
	"io"
	"net/http"
)

// GetPageHTML get page content (HTML) and then
// return it.
func GetPageHTML(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// Read body, extract HTML
	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	return string(content), nil
}
