# Use the official Golang image as a base image
FROM golang:1.24.2-bookworm

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./ 
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Command to run the application
CMD ["./main"]