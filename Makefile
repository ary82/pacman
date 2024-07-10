all: build

build:
	@echo "building go binary..."
	@go build -o main cmd/tui/main.go

run: build
	@echo "running binary..."
	@./main

clean:
	@echo "cleaning Go binary..."
	@rm main

.PHONY: all build run clean watch
