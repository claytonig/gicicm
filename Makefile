BINARY=go-icm

.PHONY: setup
setup: ## Install all the build and lint dependencies
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
	go get -u golang.org/x/tools/cmd/cover
	go mod vendor

.PHONY: testall
testall: ## Run all the tests
	make unit
	make integration

.PHONY: integration
integration: ## Run integrations tests
	docker-compose -f docker-compose.dev.yml up -d --force-recreate
	go test -cover ./... -race --tags=integration
	docker-compose -f docker-compose.dev.yml down
	docker-compose -f docker-compose.dev.yml rm -f

.PHONY: unit
unit: ## Run unit tests
	go test -cover ./... -race

.PHONY: lint
lint: ## Run all the linters
		golangci-lint run ./...

.PHONY: buildl
buildl: ## build linux binary
	GOOS=linux go build -v -o ${BINARY}


.PHONY: mod
mod: ## Run go mod for deps
	go mod vendor
