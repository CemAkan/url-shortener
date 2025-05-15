# ---------- Stage 1: Build ----------
FROM golang:1.21-alpine AS builder

# Install git, ca-certificates
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go.mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN go build -o url-shortener ./cmd/main.go

# ---------- Stage 2: Run ----------
FROM alpine:latest

# Install ca-certificates
RUN apk add --no-cache ca-certificates

# Copy built binary
COPY --from=builder /app/url-shortener /url-shortener

# Expose port
EXPOSE 3000

# Start the app
ENTRYPOINT ["/url-shortener"]