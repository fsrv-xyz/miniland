package main

import (
	"fmt"
	"log"
	"miniland/internal/cosmetic"
	"miniland/internal/filesystem"
	"miniland/internal/power"
	"miniland/internal/sysctl"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func mountfs() error {
	for _, mountpoint := range []filesystem.Mountpoint{
		{
			Source: filesystem.TMPFS,
			Target: "/tmp",
			Fstype: filesystem.TMPFS,
			Flags:  syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_RELATIME,
			Data:   "size=50M",
		},
		{
			Source: filesystem.DEVTMPFS,
			Target: "/dev",
			Fstype: filesystem.DEVTMPFS,
		},
		{
			Source: filesystem.DEVPTS,
			Target: "/dev/pts",
			Fstype: filesystem.DEVPTS,
		},
		{
			Source: filesystem.PROC,
			Target: "/proc",
			Fstype: filesystem.PROC,
		},
		{
			Source: filesystem.PROC,
			Target: "/proc",
			Fstype: filesystem.PROC,
		},
		{
			Source: filesystem.SYSFS,
			Target: "/sys",
			Fstype: filesystem.SYSFS,
		},
	} {
		err := mountpoint.Mount()
		if err != nil {
			return err
		}
	}
	return nil
}

// Function for initial booting process
func init() {
	log.Println("init: mounting filesystems")
	if err := mountfs(); err != nil {
		log.Println(err)
	}
}

func main() {
	defer func() {
		log.Println("shutting down")
		time.Sleep(20 * time.Second)
		power.Shutdown()
	}()

	err := sysctl.ApplyFile("/Config/sysctl.conf")
	if err != nil {
		log.Println(err)
	}

	cosmetic.ClearScreen()

	cmd := exec.Command("/bin/networking", "/Config/networking.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	fmt.Println("sleeping 20 seconds")
	time.Sleep(20 * time.Second)

	ief, err := net.InterfaceByName("eth0")
	if err != nil {
		fmt.Println(err)
		return
	}
	addrs, err := ief.Addrs()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", addrs)

	syscall.Chroot(".")
}
