# Этап компиляции
FROM golang:1.23 AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY ../ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o quotesbook ./cmd/main.go

# Этап сборки
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/quotesbook .

COPY ../config/config.yml ./config/config.yml

COPY internal/database/migrations ./internal/database/migrations

EXPOSE 8080

CMD ["./quotesbook"]
