
FROM golang:1.23-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./product-history-service ./cmd

FROM alpine:latest

WORKDIR /app
COPY --from=builder /build/product-history-service ./product-history-service
CMD ["/cmd/product-history-service"]
