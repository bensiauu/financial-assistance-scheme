# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/api/main.go

# Stage 2: Build a small, secure final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /root/

# Copy the binary and wait script from the builder stage
COPY --from=builder /app/main .
COPY wait-for-it.sh .

# Copy the migrations directory
COPY --from=builder /app/pkg/db/migrations/ ./pkg/db/migrations/

# Expose the application's port
EXPOSE 8080

# Command to run the app with the wait script
CMD ["./wait-for-it.sh", "db:5432", "--", "./main"]
