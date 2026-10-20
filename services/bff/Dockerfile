# syntax=docker/dockerfile:1

# ==============================================================================
# Stage 1: Build React Frontend
# ==============================================================================
FROM node:20-alpine AS frontend
WORKDIR /app

# Copy package files for the frontend
COPY services/frontend/package*.json ./
RUN npm ci

# Copy the rest of the frontend source
COPY services/frontend/ ./
RUN npm run build

# ==============================================================================
# Stage 2: Build Go BFF
# ==============================================================================
FROM golang:1.27.1-alpine AS builder
WORKDIR /app

COPY . .

# Copy the built React dist from the frontend stage into the Go embed directory
# The embed directive in static.go is //go:embed all:web/dist
COPY --from=frontend /app/dist ./services/bff/web/dist

WORKDIR /app/services/bff
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/bin/bff ./cmd/main.go

# ==============================================================================
# Stage 3: Runtime
# ==============================================================================
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/bff /app/bff
EXPOSE 8080
ENTRYPOINT ["/app/bff"]