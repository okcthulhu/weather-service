package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func WeatherHandler(weatherService WeatherService) echo.HandlerFunc {
	return func(c echo.Context) error {

		lat := c.QueryParam("lat")
		lon := c.QueryParam("lon")

		// in a production environment, we'd likely
		// check for more incorrect inputs with regex perhaps
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
	}
}
