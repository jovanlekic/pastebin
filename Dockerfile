# Dockerfile
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go dependencies
RUN go mod download

# Copy the rest of your application code
COPY . .

# Build your Go application
RUN go build -o main .

# Define the entry point command
CMD ["./main"]
