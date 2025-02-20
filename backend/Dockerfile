# Use a specific Go version that's compatible with your dependencies
FROM golang:1.22

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Update go.mod to use compatible versions
RUN go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"] 