# Use specific Alpine-based Go image
FROM golang:1.22.1-alpine3.19

# Set working directory
WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy only go.mod first
COPY go.mod ./

# Clean and verify modules
RUN go mod download && go mod verify

# Copy the rest of the code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"] 