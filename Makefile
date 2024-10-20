BINARY_NAME := main
SRC_DIR := .
BUILD_DIR := build

GO := go

.PHONY: all
all: run-all 

.PHONY: build
run-all:
	make build
	make run
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BINARY_NAME) $(SRC_DIR)
	@echo "Success to build\n"

run: 
	@echo "Execute go binary, $(BINARY_NAME)..."
	./$(BINARY_NAME)
	@echo "Success to execute!\n"

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	@echo "Success to clean up!\n"

.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  make            - Build the project"
	@echo "  make clean     - Remove the built binary"
	@echo "  make help      - Display this help message"
