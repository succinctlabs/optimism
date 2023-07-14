test:
	go test -v ./...

test-erigon:
	go test -v --l2.client=../op-erigon/runner/runner -timeout 10m

lint:
	golangci-lint run -E goimports,sqlclosecheck,bodyclose,asciicheck,misspell,errorlint -e "errors.As" -e "errors.Is"

.PHONY: \
	test \
	test-erigon \
	lint
