.PHONY: all test build docker proto

all: build test

build:
	go mod download
	go build -o ./bin/chat_service

test:
	go test ./...

docker:
	sudo docker compose up -d --build
	sudo docker logs -tf chat_service

rebuild:
	go mod tidy
	make build

format:
	go fmt ./...
	
proto:
	protoc --go_out=. --go_opt=paths=import \
    --go-grpc_out=. --go-grpc_opt=paths=import \
    proto/main.proto
	cd api && GOPROXY=direct go mod init github.com/gremislaw/chat_service/api
	cd api && GOPROXY=direct go mod tidy

clean:
	rm -rf ./bin
	rm -rf ./data