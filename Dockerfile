# Build stage
FROM golang:1.25-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod tidy

# Copy source code
COPY . .

# Build the application
RUN go build -o main .
# Final stage
FROM alpine:3.20

# System group && system user as a part of the system group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /home/appuser

# Copy the Go binary from builder
COPY --from=builder --chown=appuser:appgroup /app/main /home/appuser/main

# Create the Nyooom directory and make it world-writable to handle volume mounts
RUN mkdir -p /home/appuser/Nyooom && \
    chmod 777 /home/appuser/Nyooom

# Switch to non-root user
USER appuser

# Run the binary
CMD ["./main"]
