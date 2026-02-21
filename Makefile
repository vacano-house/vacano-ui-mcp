.PHONY: build run clean test vendor docker-build docker-run help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build         Build binary to ./bin/server"
	@echo "  run           Run without building (go run)"
	@echo "  test          Run tests"
	@echo "  vendor        Download dependencies to vendor/"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-run    Run Docker container"
	@echo "  clean         Remove ./bin directory"
	@echo "  help          Show this help"

build:
	go build -o ./bin/server ./cmd/server

run:
	go run ./cmd/server

clean:
	rm -rf ./bin

test:
	go test ./...

vendor:
	go mod tidy
	go mod vendor

docker-build:
	docker build -t vacano-ui-mcp .

docker-run:
	docker run --env-file .env -p 3000:3000 vacano-ui-mcp
