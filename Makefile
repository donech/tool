PROJECT:=tool

.PHONY: test
test:
	go test ./...
.PHONY: lint
lint:
	golangci-lint run --disable=typecheck ./...