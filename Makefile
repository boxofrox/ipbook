.PHONY: all install build

VERSION := 0.0.1
GIT_COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date -u +%Y%m%d.%H%M%S.%3N)
PKG := github.com/boxofrox/ipbook
IPBOOKD := $(PKG)/bin/ipbookd

LIBS := $(addprefix $(PKG)/lib/,pool protocol registry)
BINS := $(addprefix $(PKG)/bin/,ipbookd)

LIB_LDFLAGS =
BIN_LDFLAGS = -ldflags "-X main.VERSION $(VERSION) -X main.GIT_COMMIT $(GIT_COMMIT) -X main.BUILD_DATE $(BUILD_DATE)"

all: install

install:
	go install $(LIB_LDFLAGS) $(LIBS)
	go install $(BIN_LDFLAGS) $(BINS)

build:
	go build $(LIB_LDFLAGS) $(LIBS)
	go build $(BIN_LDFLAGS) $(BINS)

test:
	go test $(VERBOSE) ${LIB_LDFLAGS} $(LIBS)

ipbookd:
	go run $(BIN_LDFLAGS) bin/ipbookd/main.go
