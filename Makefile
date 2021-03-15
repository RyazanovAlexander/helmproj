BINDIR       := $(CURDIR)/bin
INSTALL_PATH ?= /usr/local/bin
BINNAME      ?= helmproj
BUILDDIR     ?= build
PROFILEFILE  ?= profile
TMPNAME      := $(CURDIR)/tmp
TESTTMPNAME  := $(CURDIR)/cmd/testdata/tmp
BUILDTIME    := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

# git
LASTTAG     := $(shell git tag --sort=committerdate | tail -1)
GITSHORTSHA := $(shell git rev-parse --short HEAD)

# docker option
DTAG   ?= $(LASTTAG)
DFNAME ?= Dockerfile
DRNAME ?= docker.io/aryazanov/helmproj

# example option
PROJFNAME    ?= project.yaml
EXAMPNAME    ?= examples
RENDRNAME    ?= rendered
S1NAME       ?= frontend
S2NAME       ?= backend
NSNAME       ?= example

# go option
PKG        := ./...
TESTS      := .
TESTFLAGS  :=
TAGS       :=

GOLDFLAGS += -X github.com/RyazanovAlexander/helmproj/v1/internal/version.Version=$(LASTTAG)
GOLDFLAGS += -X github.com/RyazanovAlexander/helmproj/v1/internal/version.GitShortSHA=$(GITSHORTSHA)
GOLDFLAGS += -X github.com/RyazanovAlexander/helmproj/v1/internal/version.Buildtime=$(BUILDTIME)
GOLDFLAGS += -w
GOLDFLAGS += -s
GOFLAGS   = -ldflags '$(GOLDFLAGS)'

GOOS   := linux
GOARCH := amd64

# Rebuild the buinary if any of these files change
SRC := $(shell find . -type f -name "*.go" -print) go.mod go.sum

# ------------------------------------------------------------------------------
#  run

run: build
	$(BINDIR)/$(BINNAME)

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -o $(BINDIR)/$(BINNAME) .

# ------------------------------------------------------------------------------
#  install

.PHONY: install
install: build
	@install "$(BINDIR)/$(BINNAME)" "$(INSTALL_PATH)/$(BINNAME)"

# ------------------------------------------------------------------------------
#  test

.PHONY: test
test:
	@echo
	@echo "==> Running unit tests <=="
	GO111MODULE=on go test $(GOFLAGS) -run $(TESTS) $(PKG) $(TESTFLAGS)

# ------------------------------------------------------------------------------
#  cover

.PHONY: cover
cover:
	go test -v -coverpkg=./... -coverprofile=profile ./...
	go tool cover -html=profile

# ------------------------------------------------------------------------------
#  clean

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)'
	@rm -rf '$(EXAMPNAME)/$(RENDRNAME)'
	@rm -rf '$(TMPNAME)'
	@rm -rf '$(TESTTMPNAME)'
	@rm -rf '$(PROFILEFILE)'

# ------------------------------------------------------------------------------
#  example

.PHONY: example
example:
	@helmproj -f '$(EXAMPNAME)/$(PROJFNAME)'
	@skaffold run

# ------------------------------------------------------------------------------
#  example_clear

.PHONY: example_clear
example_clear: clean
	@helm uninstall $(S1NAME) -n $(NSNAME)
	@helm uninstall $(S2NAME) -n $(NSNAME)
	@kubectl delete ns $(NSNAME)

# ------------------------------------------------------------------------------
#  container

.PHONY: container
container:
	@docker build --build-arg LDFLAGS="$(GOLDFLAGS)" --build-arg GOOS=$(GOOS) --build-arg GOARCH=$(GOARCH) -t $(DRNAME):$(DTAG) -f ./$(BUILDDIR)/$(DFNAME) .
	@docker push $(DRNAME):$(DTAG)