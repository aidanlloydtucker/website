HTTP_PORT = 8080
HTTPS_PORT = ""
CERTFILE = ""
KEYFILE = ""

# Setup the -ldflags option for go build here, interpolate the variable values

.PHONY: build clean

build:
	go build

deps:
	go get -d ./...

install: deps
	go install

clean:
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -print0 | xargs -0 rm
