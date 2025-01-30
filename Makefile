.PHONY: test
test:
	go test -cover ./...

.PHONY:build
build:
	go build ./...

.PHONY: imports
imports:
	goimports -w -local "github.com/reef-pi" ./gpio ./hal ./i2c ./pwm

.PHONY: fmt
fmt:
	gofmt -w -s ./gpio ./hal ./i2c ./pwm
