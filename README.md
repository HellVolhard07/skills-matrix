## Summary

This app introduces a consolidated setup to include multiple backend services and showcase microservice architecture with services written in different languages working seamlessly using gRPC APIs and proto definitions.

## Services Added

- **Authentication Service**
  - Uses standard HTTP requests.
  - Implemented in Go

- **Broker Service**
  - Acts as a gateway to coordinate between services.
  - Implemented in Go

- **Contact Service**
  - Sends an email to the requested person.
  - Implemented in Go and uses gRPC APIs and proto definitions.

- **User Service**
  - Handles user profiles to create, update and fetch users.
  - Implemented in python and uses gRPC APIs and proto definitions.

- **Search Service**
  - Handles searching users based on skills and proficiency.
  - Implemented in JS and uses gRPC APIs and proto definitions.

- **Skills Service**
  - Handles skills associated with users to add, update and fetch skills.
  - Implemented in Go and uses gRPC APIs and proto definitions.

- **Test Frontend**
  - Simple UI to interact with and validate backend services.

## Development Tooling

- **Makefile**
  - Added for easier build and run commands across services.

- **Docker Compose**
  - Includes all services for seamless orchestration in a local environment.

- **Mock Database**
  - Each service is using its own mock database.

## To run the app
- Clone the repo.
- Run make up_build to build the required binaries and run the services in a docker container.
- Run make start to start the frontend.
