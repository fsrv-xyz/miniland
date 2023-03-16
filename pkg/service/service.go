package service

import (
	"encoding/json"
	"io"
	"os/exec"
)

type Service struct {
	Identifier    string
	Configuration Configuration

	cmd *exec.Cmd
}

func (s *Service) ReadConfiguration(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&s.Configuration)
}

type Configuration struct {
	Name      string   `json:"name" yaml:"name"`
	Command   string   `json:"command" yaml:"command"`
	Arguments []string `json:"arguments" yaml:"arguments"`
}
