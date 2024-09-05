PACKAGE = safechan
COVER_FILE ?= coverage.out
NAMESPACE = github.com/pryg/$(PACKAGE)

help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

mod: ## Manage go mod dependencies, beautify go.mod and go.sum files
	go mod tidy && go mod vendor

fmt: ## Run go fmt for the whole project
	test -z $$(for d in $$(go list -f {{.Dir}} ./...); do gofmt -e -l -w $$d/*.go; done)

imports: ## Check and fix import section by import rules
	go install golang.org/x/tools/cmd/goimports@latest
	test -z $$(for d in $$(go list -f {{.Dir}} ./...); do goimports -e -l -local $(NAMESPACE) -w $$d/*.go; done)

lint: ## Check the project with lint
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run -c .golangci.yml --fix ./...

test: mod ## Run unit tests
	go install github.com/mfridman/tparse@latest
	go clean -testcache
	go test -coverprofile=$(COVER_FILE) ./... -json | tparse -all
	go tool cover -func=$(COVER_FILE) | grep ^total

static_check: fmt imports lint ## Run static checks (fmt, imports, lint, ...) all over the project

check: static_check test ## Run static checks and test
