# --- Stage 1: Builder ---
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker's build cache.
# This step only runs if go.mod or go.sum change.
COPY go.mod go.sum ./

# Download all Go modules (dependencies)
# This caches dependencies, so subsequent builds are faster if dependencies haven't changed.
RUN go mod download

# Copy the entire source code of your application into the container
COPY . .

# Build the Go application
# CGO_ENABLED=0: Disable CGO to create a statically linked binary, crucial for 'alpine' base image.
# GOOS=linux: Specifies the target operating system.
# -a: Force rebuilding of packages that are already up-to-date.
# -installsuffix netgo: Use the pure Go Net package for static linking.
# -o skill-server: Specifies the output binary name.
# -ldflags="-s -w": Strips debugging information and symbol tables, significantly reducing binary size.
# ./skill_server/main.go: Specifies the path to your main package.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix netgo -o skill-server ./skill_server/main.go

# --- Stage 2: Runner ---
# Use a minimal base image for the final runtime image
FROM alpine:latest

# Set the working directory for the running application
WORKDIR /app

# Copy the compiled binary from the 'builder' stage to the final image
COPY --from=builder /app/skill-server .

# Expose the gRPC port.
# Ensure this matches the port your gRPC server listens on (default 50051 in our Go example).
EXPOSE 50051

# Define the command to run the application when the container starts
ENTRYPOINT ["./skill-server"]