package main

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the HTTPClient interface
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetWeather_Success(t *testing.T) {
	mockClient := new(MockHTTPClient)
	weatherClient := NewWeatherServiceClient(mockClient, "http://example.com/points/%s,%s")

	pointsResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"properties": {"forecast": "http://example.com/forecast"}}`)),
	}
	mockClient.On("Get", "http://example.com/points/35.6895,139.6917").Return(pointsResponse, nil)

	forecastResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(`{
			"properties": {
				"periods": [{
					"name": "Today",
					"temperature": 70,
					"temperatureUnit": "F",
					"shortForecast": "Partly Cloudy"
				}]
			}
		}`)),
	}
	mockClient.On("Get", "http://example.com/forecast").Return(forecastResponse, nil)

	forecast, temp, err := weatherClient.GetWeather("35.6895", "139.6917")
	assert.NoError(t, err)
	assert.Equal(t, "Partly Cloudy", forecast)
	assert.Equal(t, 70, temp)
}

func TestGetWeather_FailToGetPoints(t *testing.T) {
	mockClient := new(MockHTTPClient)
	weatherClient := NewWeatherServiceClient(mockClient, "http://example.com/points/%s,%s")

	errorResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	mockClient.On("Get", "http://example.com/points/35.6895,139.6917").Return(errorResponse, errors.New("failed to reach API"))

	forecast, temp, err := weatherClient.GetWeather("35.6895", "139.6917")
	assert.Error(t, err)
	assert.Equal(t, "", forecast)
	assert.Equal(t, 0, temp)
}

func TestGetWeather_FailToGetForecast(t *testing.T) {
	mockClient := new(MockHTTPClient)
	weatherClient := NewWeatherServiceClient(mockClient, "http://example.com/points/%s,%s")

	pointsResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"properties": {"forecast": "http://example.com/forecast"}}`)),
	}
	mockClient.On("Get", "http://example.com/points/35.6895,139.6917").Return(pointsResponse, nil)

	errorResponse := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("")),
	}
	mockClient.On("Get", "http://example.com/forecast").Return(errorResponse, nil)

	forecast, temp, err := weatherClient.GetWeather("35.6895", "139.6917")
	assert.Error(t, err)
	assert.Equal(t, "", forecast)
	assert.Equal(t, 0, temp)
}

func TestCategorizeTemperature(t *testing.T) {
	weatherClient := NewWeatherServiceClient(nil, "")

	assert.Equal(t, "hot", weatherClient.CategorizeTemperature(90))
	assert.Equal(t, "cold", weatherClient.CategorizeTemperature(50))
	assert.Equal(t, "moderate", weatherClient.CategorizeTemperature(70))
}
