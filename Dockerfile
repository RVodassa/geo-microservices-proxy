FROM golang:1.23.3-alpine AS builder

WORKDIR /app/proxy-service

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники приложения в рабочую директорию
COPY . .

# Компилируем приложение
RUN go build -o proxy-service main.go

FROM alpine:latest

WORKDIR /app/proxy-service

# Копируем скомпилированное приложение из образа builder
COPY --from=builder /app/proxy-service .

# Копируем папку configs
COPY ./configs /app/proxy-service/configs

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./proxy-service"]