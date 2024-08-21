package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	Responses map[string]*resty.Response
}

func (m *MockClient) Get(url string) (*resty.Response, error) {
	return m.Responses[url], nil
}

func TestHttpClientGet(t *testing.T) {
	mockClient := &MockClient{
		Responses: map[string]*resty.Response{
			"https://example.com/test": {
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewBufferString(`{"message": "Hello, World!"}`),
					),
				},
				Request: &resty.Request{},
			},
		},
	}

	resp, err := mockClient.Get("https://example.com/test")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
	body, err := io.ReadAll(resp.RawResponse.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"message": "Hello, World!"}`, string(body))
}
