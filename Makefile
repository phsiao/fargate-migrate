
all: tidy fmt vet test

fmt:
	@go fmt ./...
vet:
	@go vet ./...

test:
	go test -race -cover ./...

tidy:
	@go mod tidy

.PHONY: vet fmt tidy
