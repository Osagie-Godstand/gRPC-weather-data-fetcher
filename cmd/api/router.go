package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"time"

	pb "github.com/Osagie-Godstand/gRPC-weather-data-fetcher/api/v1"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func setupRouter(grpcServerAddress string, tlsConfig *tls.Config) *mux.Router {
	router := mux.NewRouter()

	// Defining an HTTP handler to handle weather requests with city and country.
	router.HandleFunc("/weather/{city}/{country}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		city := vars["city"]
		country := vars["country"]

		conn, err := grpc.Dial(
			grpcServerAddress,
			grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := pb.NewWeatherServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		req := &pb.WeatherRequest{
			City:    city,
			Country: country,
		}

		weather, err := client.GetWeather(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Extracting the location name from the gRPC response.
		location := weather.Location

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"Location":    location,
			"Temperature": weather.Temperature,
			"Conditions":  weather.Conditions,
		})
	})

	return router
}
