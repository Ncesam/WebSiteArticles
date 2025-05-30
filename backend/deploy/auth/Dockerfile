FROM golang:1.24-alpine as builder

# Установка необходимых пакетов
RUN apk update && apk add --no-cache \
    git \
    curl \
    unzip \
    protobuf \
    make

# Установка Go-плагинов для protobuf и swag
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Работаем в директории приложения
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download


# Копируем остальное
COPY . .

# Генерация protobuf и swagger
RUN make proto

# Сборка
RUN go build -o /auth_service ./cmd/auth/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/config/.env ./config/.env
COPY --from=builder /auth_service .


CMD ["./auth_service"]
