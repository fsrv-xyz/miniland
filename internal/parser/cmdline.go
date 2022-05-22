package parser

import (
	"io/ioutil"
	"strings"
)

func ParseCmdline() (cmdline map[string]string, err error) {
	cmdline = make(map[string]string)
	b, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		return nil, err
	}
	parts := strings.Split(strings.TrimSpace(string(b)), " ")
	for _, part := range parts {
		// separate key/value based on the first = character;
		// there may be multiple (e.g. in rd.luks.name)
		if idx := strings.IndexByte(part, '='); idx > -1 {
			cmdline[part[:idx]] = part[idx+1:]
		} else {
			cmdline[part] = ""
		}
	}
	return cmdline, nil
}
