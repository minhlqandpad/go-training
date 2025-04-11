# Build stage
FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /app/cmd/app/main /app/cmd/app/main.go

# Final stage
FROM alpine:3.21

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/cmd/app/main /app/cmd/app/main

EXPOSE 50051

CMD ["/app/cmd/app/main"]
