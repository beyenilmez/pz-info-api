# 1) Build stage
FROM golang:1.23-alpine AS builder

# Disable CGO for a static binary and enable Go modules
ENV CGO_ENABLED=0 GO111MODULE=on

# Create and set the working directory
WORKDIR /app

# Copy go.mod and go.sum first so we can cache modules
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the project files
COPY . .

# Build the binary
#   -ldflags "-s -w" strips debugging info to reduce size
RUN go build -ldflags "-s -w" -o /server ./cmd/server

# 2) Final stage
FROM scratch

# Copy the compiled binary from builder stage
COPY --from=builder /server /server

# Expose port 8080
EXPOSE 8080

# Run as non-root user for better security.
USER 65532:65532

# By default, containers run as PID 1, so we just need the app entrypoint
ENTRYPOINT ["/server"]
