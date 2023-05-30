.PHONY: build
build:
	go build -o ./bin/app ./cmd/app/main.go

.PHONY: test
test:
	go test -v ./internal/...
	go test ./... -coverprofile=cover.out.tmp
	cat cover.out.tmp | grep -v "_mock.go" > cover.out
	go tool cover -func cover.out
	go tool cover -html=cover.out -o cover.html
	rm cover.out.tmp cover.out

.PHONY: lint
lint:
	golangci-lint -c .golangci.yml run ./...
