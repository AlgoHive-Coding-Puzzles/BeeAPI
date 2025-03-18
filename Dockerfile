FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o beeapi

# Final stage with Python support
FROM python:3.10-alpine

WORKDIR /app

# Copy the binary
COPY --from=builder /app/beeapi /app/beeapi

# Create puzzles directory
RUN mkdir -p /app/puzzles && chmod -R 777 /app/puzzles

# Copy static files
COPY --from=builder /app/public /app/public

# Environment variables
ENV SERVER_NAME="Local"
ENV SERVER_DESCRIPTION="Local Dev Server"
ENV PYTHON_PATH="python"

EXPOSE 8080

CMD ["/app/beeapi"]
