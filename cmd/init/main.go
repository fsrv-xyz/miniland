package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"golang.org/x/sys/unix"

	"ref.ci/fsrvcorp/miniland/userland/internal/cosmetic"
	"ref.ci/fsrvcorp/miniland/userland/internal/filesystem"
	"ref.ci/fsrvcorp/miniland/userland/internal/metrics"
	"ref.ci/fsrvcorp/miniland/userland/internal/parser"
	"ref.ci/fsrvcorp/miniland/userland/internal/power"
	"ref.ci/fsrvcorp/miniland/userland/internal/sysctl"
	"ref.ci/fsrvcorp/miniland/userland/pkg/linux/account"
	"ref.ci/fsrvcorp/miniland/userland/pkg/service"
)

func mountfs() error {
	for _, mountpoint := range []filesystem.Mountpoint{
		{
			Source: filesystem.TMPFS,
			Target: "/tmp",
			Fstype: filesystem.TMPFS,
			Flags:  syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_RELATIME,
			Data:   "size=20M",
		},
		{
			Source: filesystem.TMPFS,
			Target: "/prometheus",
			Fstype: filesystem.TMPFS,
			Flags:  syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_RELATIME,
			Data:   "size=200M",
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

	zlog.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger().With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

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
	go metrics.ServeMetrics("[::1]:9101")

	zlog.Info().Msg("starting services")
	services, err := service.DiscoverServices()
	if err != nil {
		log.Println(err)
	}
	for serviceIndex := range services {
		zlog.Info().Msgf("starting service %s", services[serviceIndex].Configuration.Name)
		services[serviceIndex].Start()
	}
	accountState := account.NewState()
	zlog.Err(accountState.AddUser(
		account.UserWithName("root"),
		account.UserWithComment("root user"),
		account.UserWithUid(0),
		account.UserWithGid(0),
		account.UserCreateGroup(
			account.GroupWithName("root"),
			account.GroupWithGid(0),
			account.GroupWithMembers([]string{"root"}),
		),
	)).Msg("adding user")
	zlog.Err(accountState.WriteFiles()).Msg("writing files")

	shell := exec.Command("/bin/sh")
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr
	shell.Stdin = os.Stdin
	if err := shell.Run(); err != nil {
		log.Println(err)
	}

	//reader := bufio.NewReader(os.Stdin)
	//fmt.Println("Simple Shell")
	//fmt.Println("---------------------")

	//for {
	//	fmt.Print("-> ")
	//	text, _ := reader.ReadString('\n')
	//	// convert CRLF to LF
	//	text = strings.Replace(text, "\n", "", -1)

	//	if strings.Compare("hi", text) == 0 {
	//		fmt.Println("hello, Yourself")
	//	}

	//}

	select {}
}
