# Stage 1: Build the Go application
FROM golang:1.20-alpine AS build

# Install dependencies
RUN apk add --no-cache git curl

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o api-gateway .

# Stage 2: Final image with only the compiled binary
FROM alpine:latest

# Install dependencies required for the runtime (if any)
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the build stage
COPY --from=build /app/api-gateway .

# Expose the port the app runs on
EXPOSE 8080

# Run the compiled binary
CMD ["./api-gateway"]
