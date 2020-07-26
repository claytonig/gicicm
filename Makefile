BINARY=go-icm

.PHONY: setup
setup: ## Install all the build and lint dependencies
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
	go get -u golang.org/x/tools/cmd/cover
	go mod vendor

.PHONY: test
test: ## Run all the tests with coverage

	echo 'mode: atomic' > coverage.txt && \
		go test -covermode=atomic -coverprofile=coverage.txt -v -race -timeout=30s \
		$$(go list ./... | grep -Ev 'vendor|mocks')

.PHONY: lint
lint: ## Run all the linters
		golangci-lint run ./...

.PHONY: buildl
buildl: ## build linux binary
	GOOS=linux go build -v -o ${BINARY}


.PHONY: mod
mod: ## Run go mod for deps
	go mod vendor
