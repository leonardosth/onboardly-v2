# =============================================
# Stage 1: Build the Go binary
# =============================================
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# =============================================
# Stage 2: Runtime image
# =============================================
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS and tzdata for timezones
RUN apk add --no-cache ca-certificates tzdata

# Copy the compiled binary from the builder stage
COPY --from=builder /build/server .

# Copy migrations (useful for initialization scripts)
COPY --from=builder /build/migrations ./migrations

# Expose the API port
EXPOSE 8080

# Run the server
ENTRYPOINT ["/app/server"]
