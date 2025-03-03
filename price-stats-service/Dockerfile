
FROM golang:1.23-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./price-stats-service ./cmd

FROM alpine:latest

WORKDIR /app
COPY --from=builder /build/price-stats-service ./price-stats-service
CMD ["/app/price-stats-service"]
