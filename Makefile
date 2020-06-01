.PHONY: build
build:
	go build -o bin/bitmedia ./cmd/bitmedia

.PHONY: build-run
build-run:
	go build -o bin/bitmedia ./cmd/bitmedia && ./bin/bitmedia

.PHONY: build-run-migrate
build-run-migrate:
	go build -o bin/bitmedia ./cmd/bitmedia && ./bin/bitmedia -migrate=true

#.PHONY: test
#test:
#	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build
