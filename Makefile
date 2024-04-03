GIT_BRANCH       := $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT       := $(shell git rev-parse HEAD)
GIT_COMMIT_SHORT := $(shell git rev-parse --short HEAD)
GIT_VERSION      := $(shell git describe --tags --always --dirty)
GIT_TAG          := $(shell git describe --tags)
BUILD_TIME       := $(shell date -u +"%Y-%m-%dT%H:%M:%S %Z")

LDFLAGS := -ldflags '\
-X "sh2unpack/constants.GitBranch=$(GIT_BRANCH)" \
-X "sh2unpack/constants.GitCommit=$(GIT_COMMIT)" \
-X "sh2unpack/constants.GitCommitShort=$(GIT_COMMIT_SHORT)" \
-X "sh2unpack/constants.GitVersion=$(GIT_VERSION)" \
-X "sh2unpack/constants.GitTag=$(GIT_TAG)" \
-X "sh2unpack/constants.BuildTime=$(BUILD_TIME)"'
BINARY_NAME = sh2unpack
GOCMD    = go
GOGEN    = $(GOCMD) generate ./
GOFORMAT = $(GOCMD) fmt
GOBUILD  = $(GOCMD) build -trimpath
GOCLEAN  = $(GOCMD) clean
GOTEST   = $(GOCMD) test
GOGET    = $(GOCMD) get
GOBUILDEXT = $(GOBUILD) -o $(BINARY_NAME) $(LDFLAGS)

.EXPORT_ALL_VARIABLES:
	GOAMD64 = v3
.PHONY: build format clean run debug deps
build:
	$(GOBUILDEXT)
	@echo Build completed
format:
	$(GOFORMAT) ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run: build
	./$(BINARY_NAME)
debug: build
	./$(BINARY_NAME) --debug
deps:
	$(GOGET)
