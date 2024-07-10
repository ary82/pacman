all: build-tui build-ssh

build-tui:
	@echo "building tui binary..."
	@go build -o tui-bin cmd/tui/main.go

run-tui: build-tui
	@echo "running tui binary..."
	@./tui-bin

build-ssh:
	@echo "building ssh binary..."
	@go build -o ssh-bin cmd/ssh/main.go

run-ssh: build-ssh
	@echo "running ssh binary..."
	@./ssh-bin

clean:
	@echo "cleaning Go binaries..."
	@rm -f tui-bin ssh-bin
	@echo "cleaning .ssh dir..."
	@rm -rf .ssh

.PHONY: all build-tui build-ssh run-tui run-ssh clean
