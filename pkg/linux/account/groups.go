package account

import (
	"fmt"
	"strconv"
	"strings"
)

type groupEntry struct {
	Name     string
	Password string
	Gid      int
	Members  []string
}

type groupFile struct {
	Entries []groupEntry
}

func (entry groupEntry) String() string {
	if entry.Name == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s:%d:%s", entry.Name, entry.Password, entry.Gid, strings.Join(entry.Members, ","))
}

func (entry groupEntry) parse(s string) (any, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid line: %s", s)
	}
	gid, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid line: %s", s)
	}
	return &groupEntry{
		Name:     parts[0],
		Password: parts[1],
		Gid:      gid,
		Members:  strings.Split(parts[3], ","),
	}, nil
}
