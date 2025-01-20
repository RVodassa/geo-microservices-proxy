FROM golang:1.23.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники приложения в рабочую директорию
COPY . .

# Компилируем приложение
RUN go build -o proxy-service ./cmd/main.go

FROM alpine

WORKDIR /root/

# Копируем скомпилированное приложение из образа builder
COPY --from=builder /app/proxy-service .

# Открываем порт 10101
EXPOSE 8080

# Запускаем приложение
CMD ["./proxy-service"]
