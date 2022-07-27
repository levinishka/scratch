.PHONY: build
build:
	mkdir -p cmd/bin
	go build -o cmd/bin/scratch cmd/scratch/main.go
	GOOS=windows GOARCH=amd64 go build -o cmd/bin/scratch.exe cmd/scratch/main.go

.PHONY: clean
clean:
	rm -rf cmd/bin
	rm -rf logs/*
	go clean

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race ./...
