# Use the official Golang 1.23.0 image as the base
FROM golang:1.23.0-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# Build the Go application
RUN go build -o bin/petplace cmd/main.go

# Set the working directory to /bin
WORKDIR /bin

# Expose port 8000
EXPOSE 8000

# Run the Go application
CMD ["petplace"]