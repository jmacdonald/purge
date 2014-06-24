GODEPS = $(realpath ./Godeps/_workspace)
GOPATH := $(GODEPS):$(GOPATH)
PATH := $(GODEPS)/bin:$(PATH)

all:
	ginkgo -r
	go build
