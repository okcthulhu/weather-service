package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read environment variables
	apiURL := os.Getenv("WEATHER_API_URL")
	if apiURL == "" {
		log.Fatal("WEATHER_API_URL is not set in the .env file")
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize the HTTP client
	var httpClient HTTPClient
	httpClient = &DefaultHTTPClient{}

	// Create an instance of the WeatherService
	var weatherService WeatherService
	weatherService = NewWeatherServiceClient(httpClient, apiURL)

	// Route to get weather information based on latitude and longitude.
	e.GET("/weather", func(c echo.Context) error {
		lat := c.QueryParam("lat")
		lon := c.QueryParam("lon")

		if lat == "" || lon == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Latitude and Longitude are required",
			})
		}

		// Get weather data
		forecast, temperature, err := weatherService.GetWeather(lat, lon)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		// Categorize the temperature
		tempCategory := weatherService.CategorizeTemperature(temperature)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"forecast": forecast,
			"temperature": map[string]interface{}{
				"value":    temperature,
				"category": tempCategory,
			},
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
