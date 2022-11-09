
ifeq ($(OS),Windows_NT)
EXE=.exe
else
EXE=
endif

BASENAME=$(notdir $(shell pwd))
PROGRAM=$(BASENAME)$(EXE)
LAST_RELEASE=

REPO=$(shell go list | head -n 1)
IMAGE=$(BASENAME)
VERSION ?= $(shell git describe --tags --always --dirty)
DOCKER=docker

.PHONY: $(PROGRAM)

all: $(PROGRAM)

compile: $(PROGRAM)

$(PROGRAM):
	go build -ldflags="-X '$(REPO)/program.Version=${VERSION}'" -o $(PROGRAM)

install:
	go install -ldflags="-X '$(REPO)/program.Version=${VERSION}'"


image: Dockerfile
	$(DOCKER) build --build-arg PROGRAM=$(BASENAME) --build-arg VERSION=$(VERSION) --build-arg BASENAME=$(BASENAME) -t $(IMAGE) .

test:
	go test -v ./...

vet:
	go vet ./...


changelog: CHANGELOG.md

CHANGELOG.md: LAST=$(shell git tag --sort=-version:refname | tail -n -1 | head -n 1)
CHANGELOG.md: TOP=$(shell git tag --sort=-version:refname |  head -n 1)
CHANGELOG.md:
	git log --pretty="* %s (%h)" $(LAST).. >> $@

hooks: .git/hooks/pre-commit

.git/hooks/pre-commit: .pre-commit-config.yaml
	pre-commit install
	pre-commit install --hook-type commit-msg


info::
	@echo BASENAME=$(BASENAME)
	@echo PROGRAM=$(PROGRAM)
	@echo IMAGE=$(IMAGE)


tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/go-critic/go-critic/cmd/gocritic@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
