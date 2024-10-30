# Use the official Go image as the base image
FROM golang:1.20-alpine

# Install git and other dependencies
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod file into the container
COPY go.mod ./

# Download and install Go dependencies
RUN go mod download

# Install CompileDaemon for hot reloading
RUN go install github.com/githubnemo/CompileDaemon@latest

# Copy the rest of the application code into the container
COPY . .

# Expose the port the app runs on
EXPOSE 8080

# Run the app with CompileDaemon for hot reloading
CMD ["CompileDaemon", "-log-prefix=false", "-build=go build -o api-gateway .", "-command=./api-gateway"]
