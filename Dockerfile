# Use the official Golang image as base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Initialize a Go module
RUN go mod init mymodule && go mod tidy

# Explicitly set GOPATH to a clean directory
RUN export GOPATH=/go && \
    go build -o filosofo

# Command to run the executable
CMD ["./filosofo"]
