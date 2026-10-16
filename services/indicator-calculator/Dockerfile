# ==============================================================================
# 1 - Build the Go binary
# ==============================================================================
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .

WORKDIR /app/services/indicator-calculator
RUN CGO_ENABLED=1 CGO_LDFLAGS="-lm" go build -ldflags="-s -w" -o /app/bin/indicator-calculator ./cmd

# ==============================================================================
# 2 - Runtime image
# ==============================================================================
FROM alpine:3.20

RUN apk --no-cache add ca-certificates libc6-compat

WORKDIR /app
COPY --from=builder /app/bin/indicator-calculator /app/indicator-calculator

ENTRYPOINT ["/app/indicator-calculator"]