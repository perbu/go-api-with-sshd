GOFILES := $(shell find . -name '*.go' -type f -not -path "./vendor/*")

bin/api: $(GOFILES)
	go build -o bin/api
lint:
	golangci-lint run
test:
	go test -v ./...