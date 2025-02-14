# Build stage
FROM golang:1.23.6-alpine AS builder

# Set working directory
WORKDIR /app

# Install required system dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY Gin-Websocket/go.mod Gin-Websocket/go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY Gin-Websocket/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy static files and configurations
COPY --from=builder /app/statics ./statics
COPY --from=builder /app/conf ./conf

# Expose the port the app runs on
EXPOSE 3333

# Command to run the application
CMD ["./main"]
