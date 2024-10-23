# Use the official Go image as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Download and install any Go dependencies
RUN go mod tidy

# Build the Go app
RUN go build -o api-gateway .

# Expose the port the app runs on
EXPOSE 8080

# Run the Go app
CMD ["./api-gateway"]
