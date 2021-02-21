BINDIR       := $(CURDIR)/bin
INSTALL_PATH ?= /usr/local/bin
BINNAME      ?= helmproj
BUILDDIR     ?= build

# docker option
DTAG   ?= 1.0.0
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
LDFLAGS    := -w -s
GOFLAGS    :=

# Rebuild the buinary if any of these files change
SRC := $(shell find . -type f -name '*.go' -print) go.mod go.sum

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME): $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(BINNAME) .

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
#  clean

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)'
	@rm -rf '$(EXAMPNAME)/$(RENDRNAME)'

# ------------------------------------------------------------------------------
#  example-run

.PHONY: example-run
example-run:
    # TODO run helmproj
	#@helmproj -f '$(EXAMPNAME)/$(PROJFNAME)'
	@skaffold run

# ------------------------------------------------------------------------------
#  example-clear

.PHONY: example-clear
example-clear:
	@helm uninstall $(S1NAME) -n $(NSNAME)
	@helm uninstall $(S2NAME) -n $(NSNAME)
	@kubectl delete ns $(NSNAME)

# ------------------------------------------------------------------------------
#  build-push-di

.PHONY: build-push-di
build-push-di:
	@docker build -t $(DRNAME):$(DTAG) -f ./$(BUILDDIR)/$(DFNAME) .
	@docker push aryazanov/helmproj:$(DTAG)