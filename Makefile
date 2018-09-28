# Go parameters
PACKAGE  = pixie
GOPATH   = $(CURDIR)/.gopath
BASE	 = $(GOPATH)/src/github.com/owensengoku/$(PACKAGE)/cmd/$(PACKAGE)	
BINPATH  = bin

export GOPATH

GOCMD   = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST  = $(GOCMD) test
GOGET   = $(GOCMD) get
GODEP   = dep
BINARY_NAME = $(CURDIR)/$(BINPATH)/$(PACKAGE)   
BINARY_UNIX = $(BINARY_NAME).linux.amd64

Q = $(if $(filter 1,$V),,@)

REGISTRY ?= registry.hub.docker.com/owensengoku
IMAGE := $(REGISTRY)/$(PACKAGE)
GIT_REF = $(shell git rev-parse --short=8 --verify HEAD)
VERSION ?= $(GIT_REF)


$(BASE):
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

all: test build

build: $(BASE)
	$Q cd $(BASE) && $(GOBUILD) -o $(BINARY_NAME) -v

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps: $(BASE)
	$Q cd $(BASE) && $(GODEP) ensure

container:
	docker build . -t $(IMAGE):$(VERSION)

push: container
	docker push $(IMAGE):$(VERSION)
	@if git describe --tags --exact-match >/dev/null 2>&1; \
	then \
		docker tag $(IMAGE):$(VERSION) $(IMAGE):latest; \
		docker push $(IMAGE):latest; \
	fi

buildlinux: $(BASE)
	$Q cd $(BASE) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v
