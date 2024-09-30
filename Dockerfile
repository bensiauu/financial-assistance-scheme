# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS go-builder

WORKDIR /app

# Copy go.mod and go.sum files first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the Go application and build the binary
COPY . .
RUN go build -o main ./cmd/api/main.go


# Stage 3: Create the final image with both Go binary and Node.js static files
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /root/

COPY --from=go-builder /app/main .

COPY --from=go-builder /app/pkg/db/migrations/ ./pkg/db/migrations/

COPY wait-for-it.sh .
RUN chmod +x wait-for-it.sh

EXPOSE 8080

CMD ["./wait-for-it.sh", "db:5432", "--", "./main"]
