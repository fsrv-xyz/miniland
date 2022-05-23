.PHONY: all
all: build

prepare_environment:
	mkdir -p stage
	rm -v rootfs.cpio.gz || true

cleanup_environment:
	rm -r stage

initrd_binary_build:
	GOOS=linux GGO_ENABLED=0 GOFLAGS=-ldflags=-w go build -o stage/init cmd/init/*.go

networking_binary_build:
	GOOS=linux GGO_ENABLED=0 GOFLAGS=-ldflags=-w go build -o stage/bin/system/networking cmd/networking/*.go

stats_prepare:
	@$(eval TMP := $(shell mktemp -d))
	wget -qO- https://github.com/prometheus/node_exporter/releases/download/v1.3.1/node_exporter-1.3.1.linux-amd64.tar.gz | tar -xz -C $(TMP)
	wget -qO- https://github.com/prometheus/prometheus/releases/download/v2.35.0/prometheus-2.35.0.linux-amd64.tar.gz | tar -xz -C $(TMP)
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
	initrd_binary_build \
	networking_binary_build \
	initrd_assemble \
	cleanup_environment

run: build
	qemu-system-x86_64 \
		-kernel vmlinuz \
		-initrd rootfs.cpio.gz \
		-m 1024 \
		-net nic,model=virtio \
		-append "test=1d bla=fasel"

sync: build
	rsync -ra --progress -e 'ssh -p22' rootfs.cpio.gz vmlinuz root@hv.infra.backend.earth:/var/lib/libvirt/iso/