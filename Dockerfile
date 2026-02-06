# Build stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o pintu-backend .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/pintu-backend .

# Copy .env if it exists (optional)
COPY .env* ./

# Expose port
EXPOSE 8080

# Run application
CMD ["./pintu-backend"]
