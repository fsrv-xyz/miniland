package service

import (
	"os/exec"
)

type WatchdogAction func(exitError error)

type Watchdog struct {
	cmd *exec.Cmd
}

func NewWatchdog(cmd *exec.Cmd) *Watchdog {
	return &Watchdog{
		cmd: cmd,
	}
}

func (w *Watchdog) Wait(action WatchdogAction) {
	commandExitError := (w.cmd).Wait()
	action(commandExitError)
}
