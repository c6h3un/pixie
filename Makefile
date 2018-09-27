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

buildlinux: $(BASE)
	$Q cd $(BASE) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v

# Left for reference
# Cross compilation
# docker-build:
# 	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v
