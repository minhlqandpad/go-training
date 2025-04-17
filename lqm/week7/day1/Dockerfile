# Build stage
FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app

# Pre-copy/cache go.mod for pre-downloading dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .
RUN go build -v -o /app/cmd/app/main /app/cmd/app/main.go

# Final stage
FROM alpine:3.21

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/cmd/app/main /app/cmd/app/main

# Copy the database/migrations folder from the builder stage
COPY --from=builder /app/database/migrations /app/database/migrations

EXPOSE 50051

CMD ["/app/cmd/app/main"]
