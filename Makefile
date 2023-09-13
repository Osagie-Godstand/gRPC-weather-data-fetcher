# Defining the source directory for .proto files
PROTO_DIR := api/v1

# Defining the list of .proto files
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

# Defining the output directories for generated code
GO_OUT_DIR := $(PROTO_DIR)
GRPC_OUT_DIR := $(PROTO_DIR)

# Generating the list of Go output files and gRPC output files
GO_OUT := $(patsubst $(PROTO_DIR)/%.proto,$(GO_OUT_DIR)/%.pb.go,$(PROTO_FILES))
GRPC_OUT := $(patsubst $(PROTO_DIR)/%.proto,$(GRPC_OUT_DIR)/%.pb.gw.go,$(PROTO_FILES))

# Default target to build and run the app
.DEFAULT_GOAL := build-and-run

# Compiling targets
compile: $(GO_OUT) $(GRPC_OUT)

$(GO_OUT): $(PROTO_FILES)
	protoc $< --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative --proto_path=$(PROTO_DIR)

$(GRPC_OUT): $(PROTO_FILES)
	protoc $< --go-grpc_out=$(GRPC_OUT_DIR) --go-grpc_opt=paths=source_relative --proto_path=$(PROTO_DIR)

# Build and Run targets
build-and-run: build-app
	@./bin/app

build-app: $(GO_OUT) $(GRPC_OUT)
	@go build -o bin/app ./cmd/api/


# Clean target
clean:
	rm -f $(GO_OUT) $(GRPC_OUT)

.PHONY: compile build-and-run build-app clean

