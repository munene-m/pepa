.PHONY: build start-app run clean

# Variables
BINARY_NAME := pepa
BUILD_DIR := .bin
MAIN_FILE := cmd/server/main.go

build:
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

run: build
	@$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

clean:
	@rm -rf $(BUILD_DIR)/$(BINARY_NAME)