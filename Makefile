.PHONY: all
all: build

CACHE_DIR := $(shell pwd)/cache

PROMETHEUS_VERSION := 2.42.0
NODE_EXPORTER_VERSION := 1.5.0
BUSYBOX_VERSION := 1.35.0-x86_64-linux-musl


prepare_environment:
	mkdir -p stage $(CACHE_DIR)
	rm -vf rootfs.cpio.gz || true

cleanup_environment:
	rm -r stage

initrd_binary_build: export CGO_ENABLED = 0
initrd_binary_build: export GOOS = linux
initrd_binary_build: export GOARCH = amd64
initrd_binary_build:
	go build -trimpath -ldflags '-s -w' -o stage/init cmd/init/main.go

networking_binary_build: export CGO_ENABLED = 0
networking_binary_build: export GOOS = linux
networking_binary_build: export GOARCH = amd64
networking_binary_build:
	go build -trimpath -ldflags '-s -w' -o stage/bin/system/networking cmd/networking/main.go

prepare_ca_certificates:
	mkdir -p stage/etc/ssl/certs
	cp /etc/ssl/certs/ca-certificates.crt stage/etc/ssl/certs/ca-certificates.crt

stats_prepare:
	mkdir -p stage/bin/stats
	@$(eval TMP := $(shell mktemp -d))

	test -f $(CACHE_DIR)/node_exporter-$(NODE_EXPORTER_VERSION).linux-amd64.tar.gz || \
		wget -qO $(CACHE_DIR)/node_exporter-$(NODE_EXPORTER_VERSION).linux-amd64.tar.gz https://github.com/prometheus/node_exporter/releases/download/v$(NODE_EXPORTER_VERSION)/node_exporter-$(NODE_EXPORTER_VERSION).linux-amd64.tar.gz
	@tar -xz -C $(TMP) -f $(CACHE_DIR)/node_exporter-$(NODE_EXPORTER_VERSION).linux-amd64.tar.gz

	test -f $(CACHE_DIR)/prometheus-$(PROMETHEUS_VERSION).linux-amd64.tar.gz || \
		wget -qO $(CACHE_DIR)/prometheus-$(PROMETHEUS_VERSION).linux-amd64.tar.gz https://github.com/prometheus/prometheus/releases/download/v$(PROMETHEUS_VERSION)/prometheus-$(PROMETHEUS_VERSION).linux-amd64.tar.gz
	@tar -xz -C $(TMP) -f $(CACHE_DIR)/prometheus-$(PROMETHEUS_VERSION).linux-amd64.tar.gz

	@cp -av $(TMP)/**/node_exporter stage/bin/stats/
	@cp -av $(TMP)/**/prometheus stage/bin/stats/
	@rm -rf $(TMP)

initrd_static_files:
	rsync -ra ./static/ stage

initrd_assemble:
	cd ./stage; find . | cpio -vo -H newc | gzip > ../rootfs.cpio.gz; cd ..

shell_prepare:
	test -f $(CACHE_DIR)/busybox || \
		wget -qO $(CACHE_DIR)/busybox-$(BUSYBOX_VERSION) https://busybox.net/downloads/binaries/$(BUSYBOX_VERSION)/busybox
	cp $(CACHE_DIR)/busybox-$(BUSYBOX_VERSION) stage/bin/busybox && chmod +x stage/bin/busybox
	stage/bin/busybox --install stage/bin

utils_prepare:
	make -C cmd/utils/ OUTDIR=../../stage/sbin

build: \
	prepare_environment \
	initrd_static_files \
	stats_prepare \
	shell_prepare \
	utils_prepare \
	prepare_ca_certificates \
	initrd_binary_build \
	networking_binary_build \
	initrd_assemble \
	cleanup_environment