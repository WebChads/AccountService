# Используем официальный образ Go
FROM golang:1.23

WORKDIR /app
COPY . .
RUN go build -o account-service ./cmd/app/main.go
EXPOSE 8082
CMD ["./account-service"]