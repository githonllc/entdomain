GOPATH ?= /tmp/gopath
GOMODCACHE ?= /tmp/gomodcache

export GOPATH
export GOMODCACHE

.PHONY: test cover lint fmt vet check

test:
	go test -count=1 -v ./...

cover:
	go test -count=1 -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1
	@rm -f coverage.out

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w -local github.com/githonllc/entdomain .

vet:
	go vet ./...

check: fmt vet test
