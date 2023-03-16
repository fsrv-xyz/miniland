package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"golang.org/x/sys/unix"

	"miniland/internal/cosmetic"
	"miniland/internal/filesystem"
	"miniland/internal/parser"
	"miniland/internal/power"
	"miniland/internal/sysctl"
	"miniland/pkg/service"
	"miniland/pkg/web"
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

	os.Mkdir("/log", 0o766)
}

func main() {
	defer func() {
		log.Println("shutting down")
		time.Sleep(20 * time.Second)
		power.Shutdown()
	}()

	if err := unix.Sethostname([]byte("testing")); err != nil {
		panic(err)
	}

	// read sysctl configuration
	err := sysctl.ApplyFile("/etc/sysctl.conf")
	if err != nil {
		log.Println(err)
	}

	cosmetic.ClearScreen()
	cmd := exec.Command("/bin/system/networking", "/etc/networking.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	fmt.Println(parser.ParseCmdline())

	go web.Start()

	sev := service.Service{
		Identifier: "prometheus",
	}

	cfg, _ := os.Open("/etc/services/prometheus.json")

	sev.ReadConfiguration(cfg)

	sev.Start()

	time.Sleep(1000 * time.Second)

	//c := make(chan bool)
	//go service.RunServices(c)
	//<-c
}
