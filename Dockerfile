# ---------- Stage 1: Build ----------
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy assets explicitly (optional since COPY . . includes it)
COPY ./email/assets /app/email/assets

RUN go build -o url-shortener ./cmd/app/main.go

# ---------- Stage 2: Run ----------
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary
COPY --from=builder /app/url-shortener /app/url-shortener

# Copy email assets
COPY --from=builder /app/email/assets /app/email/assets

EXPOSE 3000

ENTRYPOINT ["/app/url-shortener"]