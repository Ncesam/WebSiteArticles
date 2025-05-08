# ========== CONFIG ==========
BIN_DIR    := bin
PROTO_DIR  := internal/proto
DOCS_DIR   := docs
SWAG_OUT   := $(DOCS_DIR)
PROTOC_GEN_PATH := $(shell go list -f '{{ .Dir }}' -m github.com/google/protobuf)

# список сервисов (имена директорий в cmd/)
SERVICES   := gateway auth request queue config

.PHONY: all proto build clean swag tidy lint test run

## Сборка всего по порядку
all: tidy proto swag build

## Генерация protobuf
proto:
	protoc \
		-I=$(PROTO_DIR) \
		-I=third_party \
		--go_out=. \
		--go-grpc_out=. \
		$(PROTO_DIR)/*.proto \

## Шаблон для сборки одного сервиса
# cmd/<name>/main.go → bin/<name>
$(BIN_DIR)/%: cmd/%/main.go
	@mkdir -p $(BIN_DIR)
	go build -o $@ $<

## Сборка всех сервисов
build: $(addprefix $(BIN_DIR)/, $(SERVICES))

## Запуск (по умолчанию запускаем auth)
run: build
	./$(BIN_DIR)/auth

## Очистка артефактов
clean:
	rm -rf $(BIN_DIR)/*
	find . -name "*.pb.go" -type f -delete
	rm -rf $(SWAG_OUT)/*.json $(SWAG_OUT)/*.yaml $(SWAG_OUT)/docs.go
	rm -rf logs

## Генерация Swagger (swaggo)
swag:
	swag init 

## Обновление зависимостей
tidy:
	go mod tidy

## Линтинг
lint:
	golangci-lint run ./...

## Тесты
test:
	go test ./... -v
