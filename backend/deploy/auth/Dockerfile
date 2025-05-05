FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN make proto
RUN make swag

RUN go build -o /auth_service ./cmd/auth/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /auth_service .

EXPOSE 8080

CMD ["./auth_service"]
