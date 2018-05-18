
.PHONEY: all setup test cover

all: setup cover

setup:
		go get golang.org/x/tools/cmd/cover
		go get google.golang.org/appengine/urlfetch
		go get ./...

test:
		go test -v ./...

cover:
		go test -coverprofile=coverage.txt ./...
