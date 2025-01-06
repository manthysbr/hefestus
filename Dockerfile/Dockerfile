# Build stage
FROM golang:1.21-alpine AS builder

# Add necessary build tools
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hefestus ./cmd/server/main.go

# Final stage
FROM alpine:3.19

# Add necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata \
    && adduser -D appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/hefestus .
COPY --from=builder /app/.env.example .env

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./hefestus"]