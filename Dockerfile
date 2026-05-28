FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN go build -o api ./cmd/server

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
