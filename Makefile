.PHONY: test install ci-tests ci-linter go-untested

GOPATH_DIR=`go env GOPATH`

go-untested:
	go build -o go-untested ./cmd/go-untested

install:
	go install ./cmd/go-untested

test:
	go test -v ./...

ci-linter:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.30.0
	@$(GOPATH_DIR)/bin/golangci-lint run -v
