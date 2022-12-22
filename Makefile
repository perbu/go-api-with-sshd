build:
	go build -o bin/api
lint:
	golangci-lint run
test:
	go test -v ./...