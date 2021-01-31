default:
	@go build

.PHONY: clean
clean:
	@rm monkey-lang

.PHONY: lint
lint:
	@golangci-lint run