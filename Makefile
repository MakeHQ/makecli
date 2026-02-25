MODULE  := github.com/MakeHQ/makecli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "DEV")
DATE    := $(shell date -u +%Y-%m-%d)

LDFLAGS := -s -w \
	-X $(MODULE)/internal/build.Version=$(VERSION) \
	-X $(MODULE)/internal/build.Date=$(DATE)

.PHONY: build clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/makecli .

clean:
	rm -rf bin/
