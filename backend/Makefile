# ========== CONFIG ==========
APP_NAME=gateway
BIN_DIR=bin
PROTO_DIR=./internal/proto
MAIN_PATH=cmd/gateway/main.go
DOCS_DIR=docs
SWAG_OUT=./$(DOCS_DIR)

# ========== COMMANDS ==========

.PHONY: all proto build run clean swag tidy lint test

## Сборка всех целей по порядку
all: tidy proto swag build

## Генерация protobuf файлов
proto:
	protoc -I=${PROTO_DIR} --go_out=. --go-grpc_out=. ${PROTO_DIR}/*.proto

## Сборка приложения
build:
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PATH)

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

## Просмотр переменных окружения
env:
	@echo "APP_HOST: $(APP_HOST)"
	@echo "APP_PORT: $(APP_PORT)"
