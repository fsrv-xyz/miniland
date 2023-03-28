package service

import (
	"encoding/json"
	"io"
	"os/exec"
)

type Service struct {
	Identifier    string
	Configuration Configuration

	cmd    *exec.Cmd
	logger *Logger
}

func (s *Service) ReadConfiguration(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&s.Configuration)
}

type Configuration struct {
	Name        string            `json:"name" yaml:"name"`
	Owner       Owner             `json:"owner" yaml:"owner"`
	Command     string            `json:"command" yaml:"command"`
	Arguments   []string          `json:"arguments" yaml:"arguments"`
	Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`
	RunDir      string            `json:"run_dir,omitempty" yaml:"run_dir,omitempty"`
}

type Owner struct {
	UId uint32 `json:"uid" yaml:"uid"`
	GId uint32 `json:"gid" yaml:"gid"`
}
