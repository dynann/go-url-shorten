# Start with the official Go image
FROM golang:1.23 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Copy env file
COPY .env .env

# Build the Go binary
RUN go build -o main .

# Use a small image for production
FROM alpine:3.20

WORKDIR /app

# Install certificates for HTTPS if needed
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from the builder
COPY --from=builder /app/main .

# Expose the port Echo runs on (default: 1323)
EXPOSE 1323

# Run the app
CMD ["./main"]
