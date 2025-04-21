PROTO_DIR=internal/proto

proto:
	protoc -I. --go_out=. --go-grpc_out=. ${PROTO_DIR}/*.proto

build:
	go build -o bin/gateway ./cmd/gateway

clean:
	rm -rf bin/*

docs:
	swag init

run gateway:
	bin/gateway