.PHONY: build run clean test

# Build the application
build:
	go build -o bin/goflow cmd/main.go

# Run the application with a sample workflow
run: build
	./bin/goflow run -file examples/order_process.json

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Install the application
install: build
	cp bin/goflow /usr/local/bin/

# Default target
all: build 