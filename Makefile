GOLANGCI_VERSION = 1.32.0

ci: lint test
.PHONY: ci

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

lint: bin/golangci-lint
	@echo "--- lint all the things"
	@bin/golangci-lint run
.PHONY: lint

test:
	@echo "--- test all the things"
	@go test -coverprofile=coverage.txt ./...
	@go tool cover -func=coverage.txt
.PHONY: test
