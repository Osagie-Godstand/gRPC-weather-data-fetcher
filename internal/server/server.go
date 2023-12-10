package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	pb "github.com/Osagie-Godstand/gRPC-weather-data-fetcher/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WeatherServer struct {
	pb.UnimplementedWeatherServiceServer
}

// mustEmbedUnimplementedWeatherServiceServer implements weather_v1.WeatherServiceServer.
func (s *WeatherServer) mustEmbedUnimplementedWeatherServiceServer() {
	panic("unimplemented")
}

const (
	OpenWeatherMapAPIKey     = "OPENWEATHERMAP_API_KEY"
	OpenWeatherMapAPIBaseURL = "OPENWEATHERMAP_API_BASE_URL"
)

func (s *WeatherServer) GetWeather(ctx context.Context, req *pb.WeatherRequest) (*pb.WeatherResponse, error) {
	city := req.City
	country := req.Country

	weather, err := fetchWeatherData(city, country)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Weather data not found for %s, %s", city, country)
	}

	return weather, nil
}

func fetchWeatherData(city, country string) (*pb.WeatherResponse, error) {
	apiKey := os.Getenv(OpenWeatherMapAPIKey)
	baseURL := os.Getenv(OpenWeatherMapAPIBaseURL)

	// Building the API URL with the base URL, city, country, and API key.
	url := fmt.Sprintf("%s/data/2.5/weather?q=%s,%s&appid=%s", baseURL, city, country, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherData map[string]interface{}
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, err
	}

	temperatureKelvin := weatherData["main"].(map[string]interface{})["temp"].(float64)
	temperatureCelsius := temperatureKelvin - 273.15
	conditions := weatherData["weather"].([]interface{})[0].(map[string]interface{})["description"].(string)

	location := weatherData["name"].(string)

	return &pb.WeatherResponse{
		Temperature: fmt.Sprintf("%.2f°C", temperatureCelsius),
		Conditions:  conditions,
		Location:    location,
	}, nil
}
