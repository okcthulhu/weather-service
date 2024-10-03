package main

import (
	"log"
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
	e.GET("/weather", WeatherHandler(weatherService))

	e.Logger.Fatal(e.Start(":8080"))
}
