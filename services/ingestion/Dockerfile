# ==============================================================================
# 1 - Build the Go binary
# ==============================================================================
FROM golang:1.27.2-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/services/ingestion
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/ingestion ./cmd

# ==============================================================================
# 2 - Runtime image
# ==============================================================================
FROM alpine:3.24

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/bin/ingestion /app/ingestion

ENTRYPOINT ["/app/ingestion"]