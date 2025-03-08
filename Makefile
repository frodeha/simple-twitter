build: ## Build the simple_twitter server
	@mkdir -p build
	go build -o build/simple-twitter cmd/server/main.go

setup: ## Run dependencies in docker compose
	docker compose up -d

run: build setup ## Run the server on the host machine with dependencies in docker compose
	./build/simple-twitter

.PHONY: help build
help: ## Show this help screen
	@grep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help