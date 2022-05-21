package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"miniland/internal/filesystem"
	"miniland/internal/power"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

func Mkdir(path string, mode os.FileMode) error {
	if err := os.MkdirAll(path, mode); err != nil {
		return fmt.Errorf("mkdir %v: %v", path, err)
	}
	return nil
}

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
		mountpoint.Mount()
	}
	return nil
}

// Function for initial booting process
func init() {
	log.Println("init: mounting filesystems")
	if err := mountfs(); err != nil {
		fmt.Println(err)
	}
}

var cmdline = make(map[string]string)

func parseCmdline() error {
	b, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		return err
	}
	parts := strings.Split(strings.TrimSpace(string(b)), " ")
	for _, part := range parts {
		// separate key/value based on the first = character;
		// there may be multiple (e.g. in rd.luks.name)
		if idx := strings.IndexByte(part, '='); idx > -1 {
			cmdline[part[:idx]] = part[idx+1:]
		} else {
			cmdline[part] = ""
		}
	}
	return nil
}

func localAddresses() {
	interfaces, e := net.Interfaces()
	if e != nil {
		fmt.Println(e)
	}
	for _, inter := range interfaces {
		fmt.Println("Index :", inter.Index)
		fmt.Println("Name  :", inter.Name)
		fmt.Println("HWaddr:", inter.HardwareAddr)
		fmt.Println("MTU   :", inter.MTU)
		fmt.Println("Flags :", inter.Flags)
		addrs, _ := inter.Addrs()
		for _, ipaddr := range addrs {
			fmt.Println("Addr  :", ipaddr)
		}
		fmt.Println()
	}
}

func main() {
	defer func() {
		log.Println("shutting down")
		time.Sleep(20 * time.Second)
		power.Reboot()
	}()
	time.Sleep(10 * time.Second)

	parseCmdline()
	fmt.Printf("%#v\n", cmdline)
	localAddresses()

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
	fmt.Println(addrs[0])
	fmt.Println(addrs[1])

	syscall.Chroot(".")
}
