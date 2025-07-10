package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchRandomProducts makes an HTTP GET request to the given API URL
// and returns the response body as a byte slice.
func FetchRandomProducts(apiURL string) ([]byte, error) {
	// Create a new HTTP GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	// Ensure the response body is closed after the function returns
	defer resp.Body.Close()

	// Check if the HTTP status code indicates success (2xx)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
