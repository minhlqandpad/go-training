# Build stage
FROM golang:1.24.2-alpine3.21

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /app/cmd/app/main /app/cmd/app/main.go

EXPOSE 50051

CMD ["/app/cmd/app/main"]
