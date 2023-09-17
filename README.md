# gRPC-weather-data-fetcher
An Endpoint that efficiently serves secured weather data via gRPC and HTTP by combining Protocol Buffers, the Gorilla Mux router, and TLS (Transport Layer Security) encryption.

Using a self-signed TLS certicate that is PEM encoded but suitable for development purposes only.


## Generate your own cert.pem and key.pem files with this simple command
- go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

## Automating Program Compilation with a Makefile
- To generate code from weather.proto file simply use: make compile
- To build and run target simply use: make build-and-run

## Project environment variables
- GRPC_SERVER_ADDRESS=
- HTTP_SERVER_PORT=
- OPENWEATHERMAP_API_KEY=
- OPENWEATHERMAP_API_BASE_URL=


