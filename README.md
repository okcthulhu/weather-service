# Weather Service Application

This Go application provides a simple API to retrieve weather information based on latitude and longitude. It uses the National Weather Service (NWS) API to fetch weather data, categorizes the temperature, and returns the results in JSON format.

## Features

- Fetches weather forecast data from the NWS API.
- Categorizes the temperature as "hot", "cold", or "moderate".
- Simple REST API built with Echo.

## Prerequisites

- Go 1.16+
- National Weather Service (NWS) API
- `.env` file with `WEATHER_API_URL` set.

## Running the application

Run the application using:
```go run main.go```

The application will start on port 8080. (This could be made configurable)

## API Endpoints
### GET /weather
Retrieves weather information based on latitude and longitude.

### Parameters
- lat: Latitude (required)
- lon: Longitude (required)

### Example request

```curl "http://localhost:8080/weather?lat=35.6895&lon=139.6917"```

### Example response

```
{
  "forecast": "Partly Cloudy",
  "temperature": {
    "value": 70,
    "category": "moderate"
  }
}
```

