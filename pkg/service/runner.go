package service

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	zlog "github.com/rs/zerolog/log"

	"ref.ci/fsrvcorp/miniland/userland/internal/metrics"
)

func (s *Service) DefaultWatchdogActionBuilder() WatchdogAction {
	return func(exitError error) {
		metrics.ServiceState.With(prometheus.Labels{metrics.LabelServiceIdentifier: s.Identifier}).Set(float64(metrics.ServiceStateStopped))
		metrics.ServiceRestarts.With(prometheus.Labels{metrics.LabelServiceIdentifier: s.Identifier}).Inc()
		<-time.After(5 * time.Second)

		s.Stop()

		s.Start()
	}
}

func (s *Service) Start() error {
	s.context, s.contextCancelFunc = context.WithCancel(context.Background())
	cmd := exec.Command(s.Configuration.Command, s.Configuration.Arguments...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: s.Configuration.Owner.UId, Gid: s.Configuration.Owner.GId}
	cmd.SysProcAttr.Setsid = true
	cmd.SysProcAttr.Cloneflags = syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID

	// set working directory if specified
	if s.Configuration.RunDir != "" {
		cmd.Dir = s.Configuration.RunDir
	}

	// add environment variables
	for key, value := range s.Configuration.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

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

	go s.logger.Listen(s.context, s.Identifier)

	processStartError := cmd.Start()
	if processStartError != nil {
		return processStartError
	}
	s.cmd = cmd
	zlog.Debug().
		Str("pid", fmt.Sprintf("%d", s.cmd.Process.Pid)).
		Str("service", s.Identifier).
		Msg("Started service")

	go NewWatchdog(cmd).Wait(s.DefaultWatchdogActionBuilder())

	metrics.ServiceState.With(prometheus.Labels{metrics.LabelServiceIdentifier: s.Identifier}).Set(float64(metrics.ServiceStateRunning))
	return nil
}

func (s *Service) Stop() error {
	if s.cmd == nil {
		return nil
	}
	if s.contextCancelFunc != nil {
		s.contextCancelFunc()
	}
	metrics.ServiceState.With(prometheus.Labels{metrics.LabelServiceIdentifier: s.Identifier}).Set(float64(metrics.ServiceStateStopped))
	return (s.cmd).Process.Kill()
}
