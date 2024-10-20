# Use the official Golang 1.23.0 image as the base
FROM golang:1.23.0-alpine

# Set the working directory to /app
WORKDIR /app

COPY . .
RUN go mod download


# Build the Go application
RUN go build -o bin/petplace cmd/main.go
 

# Expose port 8080
EXPOSE 8080

# Run the Go application
CMD ["./bin/petplace"]