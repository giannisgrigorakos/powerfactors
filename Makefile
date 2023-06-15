help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## Perform linting
	go vet ./...
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w `find . -name '*.go'`
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

fmt: ## Format the source code
	go fmt ./...

test-local: ## Run tests locally with race detector and test coverage
	go test ./... -race -cover

docker-build: ## Build the docker image of the application
	docker build -t powerfactors .

docker-run-background: ## Run the application into a docker container on port 3000 while port mapping to local machine's port 3000
	docker run -d -p 3000:3000 powerfactors

docker-stop-containers: ## Stop all running containers
	docker stop $$(docker ps -aq)

run-local: ## Run the application locally with localhost address and port 3000
	go run cmd/powerfactors/main.go -address=127.0.0.1 -port=3000

.PHONY: help lint fmt test-local docker-build docker-run docker-stop run-local
