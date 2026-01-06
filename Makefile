.PHONY: build run clean

BINARY_NAME=notification-system
BUILD_DIR=bin

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/server/main.go

run:
	go run cmd/server/main.go

clean:
	rm -rf $(BUILD_DIR)

test:
	go test -v ./...

fmt:
	go fmt ./...

mod:
	go mod tidy
	go mod verify

.PHONY: all
all: fmt mod build
