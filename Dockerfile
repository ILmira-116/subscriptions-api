# Stage 1: сборка
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Копируем go.mod/go.sum и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарь из папки с main.go
RUN go build -o subscriptions-api ./cmd/subscriptions-api

# Stage 2: runtime
FROM alpine:latest
WORKDIR /app

# Копируем бинарь и миграции
COPY --from=builder /app/subscriptions-api .
COPY --from=builder /app/migrations ./migrations
COPY cmd/config/.env ./    

EXPOSE 8080
CMD ["./subscriptions-api"] 