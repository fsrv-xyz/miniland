# Makefile

UTILS := $(wildcard */)
OUTDIR := /tmp/bin
GO := go

.PHONY: all build $(UTILS) clean

all: build

build: $(UTILS)

GOARGS := -ldflags '-s -w' -trimpath

$(UTILS):
	CGO_ENABLED=0 $(GO) build $(GOARGS) -o $(OUTDIR)/$(basename $(subst /,,$@)) ./$@*.go

clean:
	rm -f $(OUTDIR)/*