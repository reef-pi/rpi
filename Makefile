
.PHONY:build
build:
	go build ./...

.PHONY: test
test:
	go test -cover -race ./...

.PHONY: go-get
go-get:
