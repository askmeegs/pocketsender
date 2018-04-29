.PHONY: build container

EXECUTABLE ?= pocketsender
IMAGE ?= bin/$(EXECUTABLE)
REPO = meganokeefe/pocketsender
TAG = 0.0.1

all: build

build:
	CGO_ENABLED=0 go build --ldflags '${EXTLDFLAGS}' -o ${IMAGE} github.com/m-okeefe/pocketsender

container:
	docker run -t -w /go/src/github.com/m-okeefe/pocketsender -v `pwd`:/go/src/github.com/m-okeefe/pocketsender golang:1.10.1 make
	docker build -t $(REPO):$(TAG) .
