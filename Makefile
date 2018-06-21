test:
	go test ./integration_test

build:
	go build -o iredmail-cli main.go

.DEFAULT_GOAL := build
