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
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONEY: generate
generate:
	rm -rf ./appstore/mocks/*
	rm -rf ./playstore/mocks/*
	go generate ./...

.PHONEY: update tidy update_all
update: update_all tidy

tidy:
	GO111MODULE=on GOPRIVATE="github.com/awa/*" go mod tidy

update_all:
	GO111MODULE=on GOPRIVATE="github.com/awa/*" go get -v all