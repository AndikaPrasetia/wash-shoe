start: build
	@./bin/main

build:
	@go build -o ./bin/main ./cmd/app
