# Используем официальный образ Go для сборки
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum сначала (это позволяет кэшировать зависимости)
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарный файл
RUN go build -o go-service .

# Финальный образ (меньшего размера)
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем рабочую директорию
WORKDIR /root/

# Копируем собранный бинарник из предыдущего этапа
COPY --from=builder /app/go-service .

# Открываем порт 8080
EXPOSE 8080

# Команда для запуска сервиса
CMD ["./go-service"]