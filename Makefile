GO ?= go
export GOPATH := $(CURDIR)/_vendor:$(GOPATH)

all: build

build:
	$(GO) build -o ./bin/css_extractor main.go
