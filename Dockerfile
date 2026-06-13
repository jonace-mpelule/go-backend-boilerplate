FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server ./cmd/server

FROM debian:stable-slim

WORKDIR /app

RUN useradd --system --uid 10001 appuser

COPY --from=builder /server ./server

EXPOSE 8080

USER appuser

CMD ["./server"]
