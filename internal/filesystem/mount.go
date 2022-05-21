package filesystem

import (
	"fmt"
	"os"
	"syscall"
)

type Mountpoint struct {
	Source string
	Target string
	Fstype string
	Data   string
	Flags  uintptr
}

func (mp *Mountpoint) Mount() error {
	if err := os.MkdirAll(mp.Target, 0755); err != nil {
		return err
	}

	if err := syscall.Mount(mp.Source, mp.Target, mp.Fstype, mp.Flags, mp.Data); err != nil {
		if sce, ok := err.(syscall.Errno); ok && sce == syscall.EBUSY {
			// /sys was already mounted
		} else {
			return fmt.Errorf("%v: %v", mp.Target, err)
		}
	}

	return nil
}
