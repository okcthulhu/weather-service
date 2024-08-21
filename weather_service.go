package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// WeatherService defines the interface for interacting with the weather service.
type WeatherService interface {
	GetWeather(lat, lon string) (string, int, error)
	CategorizeTemperature(temp int) string
}

// WeatherServiceClient is the concrete implementation of WeatherService.
type WeatherServiceClient struct {
	client       HTTPClient
	apiURL       string // Base URL for /points
	forecastPath string // Forecast path suffix
}

// NewWeatherServiceClient creates a new WeatherServiceClient.
func NewWeatherServiceClient(client HTTPClient, apiURL string) *WeatherServiceClient {
	return &WeatherServiceClient{
		client:       client,
		apiURL:       apiURL,
		forecastPath: "/forecast",
	}
}

// GetWeather fetches the weather forecast and temperature for the given latitude and longitude.
func (c *WeatherServiceClient) GetWeather(lat, lon string) (string, int, error) {
	// Step 1: Get the forecast endpoint from the /points API
	pointsURL := fmt.Sprintf(c.apiURL, lat, lon)
	resp, err := c.client.Get(pointsURL)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("failed to get points data: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	var pointData struct {
		Properties struct {
			Forecast string `json:"forecast"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(body, &pointData); err != nil {
		return "", 0, err
	}

	// Step 2: Get the forecast data from the forecast endpoint
	forecastURL := pointData.Properties.Forecast
	resp, err = c.client.Get(forecastURL)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("failed to get forecast data: status code %d", resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	var forecastData struct {
		Properties struct {
			Periods []struct {
				Name            string `json:"name"`
				Temperature     int    `json:"temperature"`
				TemperatureUnit string `json:"temperatureUnit"`
				ShortForecast   string `json:"shortForecast"`
			} `json:"periods"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(body, &forecastData); err != nil {
		return "", 0, err
	}

	for _, period := range forecastData.Properties.Periods {
		if period.Name == "Today" {
			return period.ShortForecast, period.Temperature, nil
		}
	}

	return "", 0, fmt.Errorf("no forecast found for 'Today'")
}

// CategorizeTemperature categorizes the temperature into "hot", "cold", or "moderate".
func (c *WeatherServiceClient) CategorizeTemperature(temp int) string {
	switch {
	case temp >= 85:
		return "hot"
	case temp <= 60:
		return "cold"
	default:
		return "moderate"
	}
}
