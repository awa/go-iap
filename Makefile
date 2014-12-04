
.PHONEY: all setup test cover

all: setup cover

setup:
		go get code.google.com/p/go.tools/cmd/cover
		go get ./...

test:
		go test -v ./...

cover:
		go test -v -coverprofile=coverage.txt -covermode=count ./appstore

