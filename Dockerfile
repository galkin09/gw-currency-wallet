FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o gw-currency-wallet ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gw-currency-wallet .
COPY --from=builder /app/config.env .
COPY --from=builder /app/internal/storages/migrations ./migrations

EXPOSE 8080

CMD ["./gw-currency-wallet"]