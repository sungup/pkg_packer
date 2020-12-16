# Build
OS=linux
ARCH=amd64

GOENV=env GOOS=$(OS) GOARCH=$(ARCH)

TEST_FLAGS =

ifeq (1,$(ON_NETWORK))
TEST_FLAGS += -tags=on_network
endif

GOCMD=go
GOBUILD=$(GOENV) $(GOCMD) build
GOTEST=$(GOCMD) test $(TEST_FLAGS)
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOCLEAN=$(GOCMD) clean
GOLIST=$(GOCMD) list -mod=vendor

BUILD_LDFLAG="-extldflags '-static'"

PROJECT_PATH=github.com/sungup/pkg_packer

TARGET_PKG_PACKER=build/amd64/pkg-packer
SOURCE_PKG_PACKER=${PROJECT_PATH}/cmd/pkg_packer
TARGET_PACKER_SAMPLE=build/amd64/pkg-packer-sample
SOURCE_PACKER_SAMPLE=${PROJECT_PATH}/cmd/pkg_packer_sample

SOURCES=$(shell find . -name "*.go")

all: build

build: $(TARGET_PKG_PACKER) $(TARGET_PACKER_SAMPLE)

format:
	@$(GOFMT) $(shell $(GOLIST) ./... | grep -v /$(PROJECT_PATH)/)

$(TARGET_PKG_PACKER): $(SOURCES) Makefile
	@$(GOBUILD) -ldflags $(BUILD_LDFLAG) -o $(TARGET_PKG_PACKER) $(SOURCE_PKG_PACKER)

$(TARGET_PACKER_SAMPLE): $(SOURCES) Makefile
	@$(GOBUILD) -ldflags $(BUILD_LDFLAG) -o $(TARGET_PACKER_SAMPLE) $(SOURCE_PACKER_SAMPLE)

test: format
	@$(GOVET) $(shell $(GOLIST) ./... | grep -v /$(PROJECT_PATH)/)
	@$(GOTEST) $(shell $(GOLIST) ./... | grep -v /$(PROJECT_PATH)/)

gomod-refresh:
	@rm -rf go.sum go.mod vendor
	@$(GOCMD) mod init
	@$(GOCMD) mod vendor -v

clean:
	@$(GOCLEAN)
	@$(GOCLEAN) -testcache
	@rm -f $(TARGET_PKG_PACKER) $(TARGET_PACKER_SAMPLE)

