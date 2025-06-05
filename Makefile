# Makefile

.PHONY: proto build_user_service run_user_service_local run_user_service_docker test_user_service clean_docker up down

# Define variables
USER_SERVICE_DIR := user_service
GATEWAY_DIR := gateway
PROTO_FILE := $(USER_SERVICE_DIR)/user_service.proto

# Default target
all: up

# Generate protobuf files
proto:
	@echo "Generating Python protobuf and gRPC stubs..."
	python -m grpc_tools.protoc \
	    -I. \
	    -I$(USER_SERVICE_DIR) \
	    -I./proto \
	    --python_out=$(USER_SERVICE_DIR) \
	    --grpc_python_out=$(USER_SERVICE_DIR) \
	    $(PROTO_FILE)
	@echo "Applying manual fix to $(USER_SERVICE_DIR)/user_service_pb2_grpc.py..."
	# This sed command is for Linux/macOS. For Windows, you might need a different tool or manual edit.
	# It replaces 'import user_service_pb2' with 'from . import user_service_pb2'
	# Be very careful if running on Windows without a suitable 'sed' equivalent (e.g., Git Bash provides it)
	# Find more robust cross-platform ways if this causes issues.
	sed -i 's/^import user_service_pb2 as user__service__pb2/from . import user_service_pb2 as user__service__pb2/' $(USER_SERVICE_DIR)/user_service_pb2_grpc.py || true

# Build User Service Docker image
build_user_service: proto
	@echo "Building user_service Docker image..."
	docker build -t skillmetrix_user_service -f $(USER_SERVICE_DIR)/Dockerfile $(USER_SERVICE_DIR)

# Run User Service locally (without Docker)
run_user_service_local: proto
	@echo "Running User Service locally..."
	python -m $(USER_SERVICE_DIR).server

# Run User Service using Docker Compose
run_user_service_docker: build_user_service
	@echo "Running User Service via Docker Compose..."
	docker compose up --build user_service

# Run all services defined in docker-compose.yaml
up: proto # Ensure protobufs are generated before building
	@echo "Starting all services via Docker Compose..."
	docker compose up --build -d # -d for detached mode

# Stop and remove all services defined in docker-compose.yaml
down:
	@echo "Stopping and removing all services via Docker Compose..."
	docker compose down

# Run User Service tests (if you re-add them)
# test_user_service: proto
# 	@echo "Running User Service unit tests..."
# 	python -m unittest discover $(USER_SERVICE_DIR)/tests

# Clean Docker artifacts
clean_docker:
	@echo "Cleaning Docker images and containers..."
	docker stop skillmetrix_user_service || true
	docker rm skillmetrix_user_service || true
	docker rmi skillmetrix_user_service || true
	docker system prune -f --volumes # Prune unused Docker data (use with caution)