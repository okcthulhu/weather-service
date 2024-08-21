package main

import (
	"net/http"
)

// HTTPClient interface to abstract the HTTP client for easier testing.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// DefaultHTTPClient is the default implementation of HTTPClient.
type DefaultHTTPClient struct{}

// Get sends a GET request to the specified URL.
func (c *DefaultHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}
