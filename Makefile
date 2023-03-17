.PHONY: all
all: build

prepare_environment:
	mkdir -p stage
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
	wget -qO- https://github.com/prometheus/node_exporter/releases/download/v1.5.0/node_exporter-1.5.0.linux-amd64.tar.gz | tar -xz -C $(TMP)
	wget -qO- https://github.com/prometheus/prometheus/releases/download/v2.42.0/prometheus-2.42.0.linux-amd64.tar.gz | tar -xz -C $(TMP)
	@cp -av $(TMP)/**/node_exporter stage/bin/stats/
	@cp -av $(TMP)/**/prometheus stage/bin/stats/
	@rm -rf $(TMP)

initrd_static_files:
	rsync -ra ./static/ stage

initrd_assemble:
	cd ./stage; find . | cpio -vo -H newc | gzip > ../rootfs.cpio.gz; cd ..

build: \
	prepare_environment \
	initrd_static_files \
	stats_prepare \
	prepare_ca_certificates \
	initrd_binary_build \
	networking_binary_build \
	initrd_assemble \
	cleanup_environment

run: build
	qemu-system-x86_64 \
		-kernel bzImage \
		-initrd rootfs.cpio.gz \
		-m 1024 \
		-smp 2 \
		-net nic,model=e1000 \
		-append "test=1d bla=fasel"