# Step 1: Builder stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies (cached if go.mod/sum don't change)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# We build the binary for Linux target to ensure it runs in Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

# Step 2: Final runtime image
FROM alpine:latest

WORKDIR /app

# Install necessary libraries (e.g., CA-certificates for TLS)
RUN apk --no-cache add ca-certificates tzdata

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .
COPY .env .

# Expose port 4000
EXPOSE 4000

# Run the binary
CMD ["./main"]
