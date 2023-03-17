package service

import (
	"fmt"
	"os/exec"
	"time"

	zlog "github.com/rs/zerolog/log"
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
	go s.logger.Listen(s.Identifier)

	processStartError := cmd.Start()
	if processStartError != nil {
		return processStartError
	}
	s.cmd = cmd
	zlog.Debug().
		Str("pid", fmt.Sprintf("%d", s.cmd.Process.Pid)).
		Str("service", s.Identifier).
		Msg("Started service")

	go NewWatchdog(cmd).Wait(func(exitError error) {
		<-time.After(5 * time.Second)
		s.Start()
	})

	return nil
}

func (s *Service) Stop() error {
	if s.cmd == nil {
		return nil
	}
	return (s.cmd).Process.Kill()
}
