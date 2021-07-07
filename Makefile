NAME_API=library-api
VERSION=dev
OS ?= linux
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_BUILD_TIME=$(shell date '+%Y-%m-%d__%I:%M:%S%p')
DOCKER_COMPOSE_FILE=docker-compose.yml
SQLC_CONFIG_FILE=pkg/gateway/repositories/sqlc/config/sqlc.yaml

.PHONY: test
test:
	@echo "==> Running Tests"
	go test -race -v ./...

.PHONY: compile
compile: clean
	@echo "==> Go Building API"
	@env GOOS=${OS} GOARCH=amd64 go build -v -o  build/${NAME_API} \
	-ldflags "-X main.BuildGitCommit=$(GIT_COMMIT) -X main.BuildTime=$(GIT_BUILD_TIME)" ./cmd/api/


.PHONY: clean
clean:
	@echo "==> Cleaning releases"
	@GOOS=${OS} go clean -i -x ./...
	@rm -f build/${NAME_API}
	@rm -f coverage.html
	@rm -f coverage.out

.PHONY: metalint
metalint:

ifeq (, $(shell which $$(go env GOPATH)/bin/golangci-lint))
	@echo "==> installing golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	go install ./...
	go test -i ./...
endif

	$$(go env GOPATH)/bin/golangci-lint run -c ./.golangci.yml ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests"
	@richgo test -failfast -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

### remove sqlc-dev binary once new tag is released on https://github.com/kyleconroy/sqlc
.PHONY: generate
generate:
	@echo "Go Generating"
	go get github.com/matryer/moq@v0.2.1
	go get -u github.com/swaggo/swag/cmd/swag@v1.6.7
	go generate ./...
	swag init -g ./cmd/api/main.go -o ./docs/swagger
	@./pkg/gateway/repositories/sqlc/config/sqlc-dev generate -f $(SQLC_CONFIG_FILE)
	go mod tidy

.PHONY: setup-dev
setup-dev:
	@echo "Setting up dev environment"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d