package sysctl

import (
	"os"
	"strings"
)

type path []string

var BasePath = []string{"proc", "sys"}

type Property struct {
	path  path
	value string
}

func (property *Property) GetValue() string {
	return property.value
}

func (property *Property) SetValue(value string) error {
	return writeFile(property.path.filePath(), value)
}

func NewProperty(options ...PropertyOption) *Property {
	property := &Property{}
	for _, option := range options {
		option(property)
	}
	return property
}

func encodePathSlashes(pathString string) path {
	return strings.Split(pathString, "/")
}
func encodePathDots(pathString string) path {
	return strings.Split(pathString, ".")
}

func writeFile(path, value string) error {
	return os.WriteFile(path, []byte(value), 0o644)
}

func (p path) array() []string {
	return p
}

func (p path) filePath() string {
	return "/" + strings.Join(append(BasePath, p.array()...), "/")
}
