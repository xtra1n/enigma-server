.PHONY: build run test clean deps help

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server -port 8080

dev:
	go run cmd/server/main.go -port 8080

test:
	go test -v ./internal/...

clean:
	rm -rf bin/

deps:
	go mod download
	go mod tidy

help:
	@echo "Available commands:"
	@echo "  make build    - Build the server binary"
	@echo "  make run      - Build and run the server"
	@echo "  make dev      - Run with go run (development)"
	@echo "  make test     - Run all tests"
	@echo "  make clean    - Remove bin/ directory"
	@echo "  make deps     - Download and tidy dependencies"
	@echo "  make help     - Show this help"