.PHONEY: all
all: setup cover

.PHONEY: setup
setup:
	go get golang.org/x/tools/cmd/cover
	go get google.golang.org/appengine/urlfetch
	go get ./...

.PHONEY: test
test:
	go test -v ./...

.PHONEY: cover
cover:
	go test -coverprofile=coverage.txt ./...

.PHONEY: generate
generate:
	rm -rf ./appstore/mocks/*
	rm -rf ./playstore/mocks/*
	go generate ./...
