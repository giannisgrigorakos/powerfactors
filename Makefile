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

.PHONY: help lint fmt test-local
