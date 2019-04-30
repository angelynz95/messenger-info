.PHONY: all

-include .env
export $(shell sed 's/=.*//' .env)

compile: pretty
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go

dep: dep-remove dep-add dep-install

dep-add:
	@govendor add +external

dep-install:
	@govendor fetch -v +outside
	@govendor sync

dep-remove:
	@govendor remove +unused

init:
	@go get -u github.com/kardianos/govendor
	@go get -u github.com/Fs02/kamimai/cmd/kamimai

pretty:
	@go fmt ./...

run:
	@./main
