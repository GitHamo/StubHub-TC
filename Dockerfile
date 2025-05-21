FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /traffic-controller ./cmd/server

# Use a minimal alpine image for the final container
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /traffic-controller .
COPY --from=builder /app/.env .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./traffic-controller"]