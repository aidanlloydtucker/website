PORT = 8080
CERTFILE = ""
KEYFILE = ""

# Setup the -ldflags option for go build here, interpolate the variable values

LDFLAGS += -X \"main.Port=$(PORT)\"
LDFLAGS += -X \"main.CertFile=$(CERTFILE)\"
LDFLAGS += -X \"main.KeyFile=$(KEYFILE)\"

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
