# ==============================================================================
# 1 - Build the Go binary
# ==============================================================================
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/services/signal-decision
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/signal-decision ./cmd

# ==============================================================================
# 2 - Runtime image
# ==============================================================================
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/bin/signal-decision /app/signal-decision

ENTRYPOINT ["/app/signal-decision"]