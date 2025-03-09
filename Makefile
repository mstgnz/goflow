.PHONY: build run clean test test-coverage docker-build docker-run docker-compose docker-clean

# Build the application
build:
	go build -o bin/goflow cmd/main.go

# Run the application with a sample workflow
run: build
	./bin/goflow run -file examples/order_process.json

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf coverage/

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	mkdir -p coverage
	go test -v -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

# Install the application
install: build
	cp bin/goflow /usr/local/bin/

# Build Docker image
docker-build:
	docker build -t goflow:latest .

# Run Docker container
docker-run: docker-build
	docker run --rm goflow:latest

# Run with Docker Compose
docker-compose:
	docker-compose up --build

# Clean Docker resources
docker-clean:
	docker-compose down
	docker rmi goflow:latest

# Default target
all: build 