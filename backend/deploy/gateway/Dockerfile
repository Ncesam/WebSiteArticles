
FROM golang:1.21-alpine as builder

WORKDIR /app

COPY ../../go.mod ./
COPY ../../go.sum ./
RUN go mod download

COPY ../../. .

RUN make proto
RUN make swag


RUN go build -o /gateway ./cmd/gateway/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /gateway .

EXPOSE 8080

CMD ["./gateway"]
