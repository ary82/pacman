all: build

build:
	@echo "building go binary..."
	@go build -o main main.go

run: build
	@echo "running binary..."
	@./main

clean:
	@echo "cleaning Go binary..."
	@rm main

# Live Reload
watch:
	@${HOME}/go/bin/air

.PHONY: all build run clean watch
