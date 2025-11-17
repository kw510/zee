# Use an official Golang runtime as a parent image
FROM golang:1-alpine AS builder

# Set the working directory to /app
WORKDIR /app

# Copy the module files
COPY . .

# Download the go dependencies
RUN go mod download

# Build a static application binary
RUN CGO_ENABLED=0 GOOS=linux go build -o ./tmp/main

## Development stage, using air for hot reloading
FROM builder AS development
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]

## Production stage, using a static binary and scratch image
FROM scratch
COPY --from=builder /app/tmp/main /app
CMD ["/app"]