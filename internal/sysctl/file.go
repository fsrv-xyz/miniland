package sysctl

import (
	"log"
	"os"
	"strings"
)

func ApplyFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		// skip if line is a comment
		if strings.HasPrefix(line, "#") {
			continue
		}
		pair := strings.Split(line, "=")
		property := NewProperty(OptionPathFormatDots(pair[0]))
		setPropertyValueError := property.SetValue(pair[1])
		if setPropertyValueError != nil {
			log.Println(setPropertyValueError)
			continue
		}
	}
	return nil
}
