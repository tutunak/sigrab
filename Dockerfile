# Builder stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy all project files
COPY . .

# Build the Go project
RUN go build -o /sigrab cmd/sigrab/main.go

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates

# Copy the built executable from the builder stage
COPY --from=builder /sigrab /sigrab

# Set the entrypoint
ENTRYPOINT ["/sigrab"]
