.PHONY: build test vet check clean lint

build:
	go build ./...

test:
	go test ./... -count=1

vet:
	go vet ./...

check: vet test build

lint:
	golangci-lint run ./...

clean:
	rm -f claudekit claudekit-mcp

-include $(HOME)/hairglasses-studio/dotfiles/make/pipeline.mk
