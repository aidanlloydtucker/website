PORT = 8080

# Setup the -ldflags option for go build here, interpolate the variable values

LDFLAGS += -X \"main.Port=$(PORT)\"

.PHONY: build clean

build:
	go build -ldflags "$(LDFLAGS)"

install:
	go get -d ./...
	go install -ldflags "$(LDFLAGS)"

clean:
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -print0 | xargs -0 rm
