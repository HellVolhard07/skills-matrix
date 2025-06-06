<<<<<<< HEAD
.PHONY: build_user_service run_user_service_local run_user_service_docker test_user_service clean_docker \
        build_broker build_auth build_contact build_front start stop up up_build down

SHELL=cmd.exe
FRONT_END_BINARY=frontApp.exe
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
CONTACT_BINARY=contactApp

USER_SERVICE_DIR := user_service

## ====================
## Docker Compose
## ====================

up:
	@echo Starting Docker images...
	docker compose up -d
	@echo Docker images started!

up_build: build_broker build_auth build_contact build_user_service
	@echo Stopping docker images (if running...)
	docker compose down
	@echo Building (when required) and starting docker images...
	docker compose up --build -d
	@echo Docker images built and started!

down:
	@echo Stopping docker compose...
	docker compose down
	@echo Done!

clean_docker:
	@echo Cleaning Docker images and containers...
	docker stop skillmetrix_user_service || true
	docker rm skillmetrix_user_service || true
	docker rmi skillmetrix_user_service || true
	docker system prune -f --volumes

## ====================
## User Service
## ====================

build_user_service:
	@echo "Building user_service Docker image..."
	docker build -t skillmetrix_user_service -f $(USER_SERVICE_DIR)/user_service.dockerfile $(USER_SERVICE_DIR)

run_user_service:
	@echo "Running User Service container..."
	docker run --rm -d --name skillmetrix_user_service -p 50052:50052 skillmetrix_user_service

# test_user_service:
# 	@echo "Running User Service unit tests..."
# 	python -m unittest discover $(USER_SERVICE_DIR)/tests

## ====================
## Go Microservices
## ====================

build_broker:
	@echo Building broker binary...
	chdir .\broker-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${BROKER_BINARY} ./cmd/api
	@echo Done!

build_auth:
	@echo Building auth binary...
	chdir .\authentication-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${AUTH_BINARY} ./cmd/api
	@echo Done!

build_contact:
	@echo Building contact binary...
	chdir .\contact-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${CONTACT_BINARY} ./cmd/api
	@echo Done!

## ====================
## Front End
## ====================

build_front:
	@echo Building front end binary...
	chdir .\front-end && set CGO_ENABLED=0&& set GOOS=windows&& go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo Done!

start: build_front
	@echo Starting front end...
	chdir .\front-end && start /B ${FRONT_END_BINARY} &

stop:
	@echo Stopping front end...
	@taskkill /IM "${FRONT_END_BINARY}" /F
	@echo "Stopped front end!"
=======
>>>>>>> 86a1e051fa56b74af45487995c66a00cc3fab704
