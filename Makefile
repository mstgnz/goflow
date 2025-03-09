.PHONY: build run clean test test-coverage

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

# Default target
all: build 