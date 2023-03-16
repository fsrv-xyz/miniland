package service

import (
	"fmt"
	"os/exec"
)

func (s *Service) Start() error {
	cmd := exec.Command(s.Configuration.Command, s.Configuration.Arguments...)
	s.logger = &Logger{}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	s.logger.Bind(stdout, stderr)
	go s.logger.Listen()

	processStartError := cmd.Start()
	if processStartError != nil {
		return processStartError
	}
	s.cmd = cmd
	fmt.Printf("Started %s, PID: %d\n", s.Configuration.Name, s.cmd.Process.Pid)
	return nil
}

func (s *Service) Stop() error {
	if s.cmd == nil {
		return nil
	}
	return (s.cmd).Process.Kill()
}
