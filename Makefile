.PHONY: setup

setup:
	go mod download
	go run github.com/playwright-community/playwright-go/cmd/playwright install
	@echo "Setup complete" 