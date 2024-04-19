package external

import (
	"log"
	"net/http"
)

func PrepareHTTPRequest(url string) (*http.Request, error) {
	log.Printf("Preparing request to %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("user-agent", "Go/Daily-Dashboard")

	return req, nil
}
