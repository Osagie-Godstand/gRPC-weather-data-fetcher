package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	pb "github.com/Osagie-Godstand/gRPC-weather-data-fetcher/api/v1"
	"github.com/Osagie-Godstand/gRPC-weather-data-fetcher/internal/server"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	grpcServerAddress := os.Getenv("GRPC_SERVER_ADDRESS")
	httpServerPort := os.Getenv("HTTP_SERVER_PORT")

	// Loading TLS certificates
	certFile := "tls/cert.pem"
	keyFile := "tls/key.pem"

	// Loading the server's certificate
	cert, err := os.ReadFile(certFile)
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}

	// Creating a certificate pool and add the server's certificate
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("Failed to append server certificate to the certificate pool")
	}

	// Creating a TLS configuration with the certificate pool
	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	// Creating a new gRPC server with TLS credentials
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load TLS certificates: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterWeatherServiceServer(grpcServer, &server.WeatherServer{})

	router := mux.NewRouter()

	go func() {
		lis, err := net.Listen("tcp", grpcServerAddress)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		defer lis.Close()

		fmt.Printf("gRPC server is running on %s...\n", grpcServerAddress)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

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

	fmt.Printf("HTTP server is running on port %s...\n", httpServerPort)
	http.ListenAndServe(httpServerPort, router)
}
