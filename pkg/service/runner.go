package service

import (
	"fmt"
	"os"
	"os/exec"
)

func (s *Service) Start() error {
	cmd := exec.Cmd{
		Path: s.Configuration.Command,
		Args: s.Configuration.Arguments,
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	processStartError := cmd.Start()
	if processStartError != nil {
		return processStartError
	}
	s.cmd = &cmd
	fmt.Printf("Started %s, PID: %d\n", s.Configuration.Name, s.cmd.Process.Pid)
	return nil
}

func (s *Service) Stop() error {
	if s.cmd == nil {
		return nil
	}
	return (s.cmd).Process.Kill()
}
