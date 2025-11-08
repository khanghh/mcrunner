.DEFAULT_GOAL := mcrunner

BUILD_DIR=$(CURDIR)/build/bin
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DATE=$(shell git show -s --format=%cI HEAD)
GIT_TAG=$(shell git describe --tags --always --dirty)

LDFLAGS=-ldflags "-w -s -X 'main.gitCommit=$(GIT_COMMIT)' -X 'main.gitDate=$(GIT_DATE)' -X 'main.gitTag=$(GIT_TAG)'"

mcrunner:
	@echo "Building target: $@" 
	go build $(LDFLAGS) -o $(BUILD_DIR)/$@ $(CURDIR)/main.go
	@echo "Done building."

build-docker:
	@echo "Building Docker image: mcrunner:latest" 
	docker build --rm --progress=plain \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg GIT_DATE=$(GIT_DATE) \
		--build-arg GIT_TAG=$(GIT_TAG) \
		-t mcrunner:latest \
		-f ./docker/Dockerfile .
	@echo "Done building."

clean:
	@rm -rf $(BUILD_DIR)/*

all: mcrunner
