
.PHONY: test
test:
	go test -cover -race ./...

.PHONY:build
build:
	go build ./...


.PHONY: go-get
go-get:
