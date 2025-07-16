# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod .
COPY go.sum .

# Download Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# CGO_ENABLED=0 disables CGO, making the binary statically linked and suitable for a minimal base image
# -o app specifies the output binary name
# ./... builds all packages in the current module
RUN CGO_ENABLED=0 go build -o /app/crowdstrike-rtr-app ./main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/crowdstrike-rtr-app .

# Copy the .env file (important for runtime configuration)
# Ensure .env is in the same directory as your Dockerfile when building
COPY .env .

# Command to run the application when the container starts
CMD ["./crowdstrike-rtr-app"]
