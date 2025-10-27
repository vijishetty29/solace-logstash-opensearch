# Stage 1: Build the Go binary
FROM golang:1.25.2-alpine  AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/connector .

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the 'builder' stage
COPY --from=builder /app/connector .

# Run the connector
CMD ["./connector"]