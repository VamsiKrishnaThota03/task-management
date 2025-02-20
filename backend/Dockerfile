# Use specific Alpine-based Go image
FROM golang:1.22.0-alpine3.19

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Install build dependencies
RUN apk add --no-cache git

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"] 