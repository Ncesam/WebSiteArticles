# ========== CONFIG ==========
GATEWAY_NAME=gateway
AUTH_NAME=auth_service
BIN_DIR=bin
PROTO_DIR=./internal/proto
GATEWAY_PATH=cmd/gateway/main.go
AUTH_PATH=cmd/auth/main.go
DOCS_DIR=docs
SWAG_OUT=./$(DOCS_DIR)

# ========== COMMANDS ==========

.PHONY: all proto build run clean swag tidy lint test

## Сборка всех целей по порядку
all: tidy proto swag build

## Генерация protobuf файлов
proto:
	protoc -I ./internal/proto --go_out=. --go-grpc_out=. ./internal/proto/*.proto

## Сборка приложения
build:
	go build -o $(BIN_DIR)/$(GATEWAY_NAME) $(GATEWAY_PATH)
	go build -o $(BIN_DIR)/$(AUTH_NAME) $(AUTH_PATH)
## Запуск приложения
run: build
	./$(BIN_DIR)/$(APP_NAME)

## Удаление сгенерированных бинарников и временных файлов
clean:
	rm -rf $(BIN_DIR)/*
	find ./generated/ -name "*.pb.go" -type f -delete
	rm -rf $(SWAG_OUT)/*.json $(SWAG_OUT)/*.yaml $(SWAG_OUT)/docs.go
	rm -rf logs/logs.log

## Генерация Swagger документации
swag:
	swag init 

## Обновить зависимости
tidy:
	go mod tidy

## Проверка стиля кода
lint:
	golangci-lint run ./...

## Запуск тестов
test:
	go test ./... -v


