build:
	@go build -o /Users/abhishek/Documents/projects/goland/gobank/bin/gobank

run: build
	@bin/gobank

test:
	@go test ./...