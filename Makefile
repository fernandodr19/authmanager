NAME_API=library-api
VERSION=dev
OS ?= linux
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_BUILD_TIME=$(shell date '+%Y-%m-%d__%I:%M:%S%p')
DOCKER_COMPOSE_FILE=dev/docker-compose.yml

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
	@rm coverage.html
	@rm coverage.out

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

.PHONY: generate
generate:
	@echo "Go Generating"
	go get github.com/matryer/moq@v0.1.4
	go get -u github.com/swaggo/swag/cmd/swag@v1.6.7
	go generate ./...
	swag init -g ./cmd/api/main.go -o ./docs/swagger #here because this needs to be called on root folder
	go mod tidy

.PHONY: setup-dev
setup-dev:
	@echo "Setting up dev enviroment"
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	docker-compose -f $(DOCKER_COMPOSE_FILE) up