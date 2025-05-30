
FROM golang:1.24-alpine as builder

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


RUN make proto


RUN go build -o /config_service ./cmd/config/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/config/.env ./config/.env
COPY --from=builder /config_service .

CMD ["./config_service"]
