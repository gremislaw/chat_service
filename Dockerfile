#Docker Pipeline
FROM golang:1.23.2-alpine3.20 as builder
WORKDIR /app
COPY . /app
RUN go mod download && \
    go build -o ./bin/chat_service

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/bin/chat_service ./bin/chat_service
EXPOSE 9000
ENTRYPOINT ["./bin/chat_service"]
