.PHONY: build test vet check clean

build:
	go build ./...

test:
	go test ./... -count=1

vet:
	go vet ./...

check: vet test build

clean:
	rm -f claudekit claudekit-mcp
