.PHONY: all
all: build

prepare_environment:
	mkdir -p stage/bin
	rm -v rootfs.cpio.gz || true

cleanup_environment:
	rm -r stage

initrd_binary_build:
	GOOS=linux GGO_ENABLED=0 GOFLAGS=-ldflags=-w go build -o stage/init cmd/init/*.go

networking_binary_build:
	GOOS=linux GGO_ENABLED=0 GOFLAGS=-ldflags=-w go build -o stage/bin/networking cmd/networking/*.go

initrd_static_files:
	rsync -rav ./static/ stage

initrd_assemble:
	cd ./stage; find . | cpio -o -H newc | gzip > ../rootfs.cpio.gz; cd ..

build: \
	prepare_environment \
	initrd_static_files \
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
	scp -P22 -6 rootfs.cpio.gz root@hv.infra.backend.earth:/var/lib/libvirt/iso/
	scp -P22 -6 vmlinuz root@hv.infra.backend.earth:/var/lib/libvirt/iso/