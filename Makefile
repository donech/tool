PROJECT:=tool

.PHONY: test
test:
	go test ./...
.PHONY: lint
lint:
	golangci-lint --go=1.17 run  ./...
