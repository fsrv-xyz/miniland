package power

import "syscall"

func Reboot() {
	const LINUX_REBOOT_MAGIC1 uintptr = 0xfee1dead
	const LINUX_REBOOT_MAGIC2 uintptr = 672274793
	const LINUX_REBOOT_CMD_RESTART uintptr = 0x1234567
	syscall.Syscall(syscall.SYS_REBOOT,
		LINUX_REBOOT_MAGIC1,
		LINUX_REBOOT_MAGIC2,
		LINUX_REBOOT_CMD_RESTART)

}

// Shutdown - shutdown the system
func Shutdown() {
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}
