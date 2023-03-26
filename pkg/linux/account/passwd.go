package account

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type passwdEntry struct {
	Username string
	Password string
	Uid      int
	Gid      int
	Gecos    string
	Home     string
	Shell    string
}

type passwdFile struct {
	Entries []passwdEntry
}

func (entry passwdEntry) String() string {
	if entry.Username == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s:%d:%d:%s:%s:%s", entry.Username, entry.Password, entry.Uid, entry.Gid, entry.Gecos, entry.Home, entry.Shell)
}

func (entry passwdEntry) parse(s string) (any, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 7 {
		return nil, errors.New("invalid account entry")
	}
	p := &passwdEntry{}
	p.Username = parts[0]
	p.Password = parts[1]
	p.Uid, _ = strconv.Atoi(parts[2])
	p.Gid, _ = strconv.Atoi(parts[3])
	p.Gecos = parts[4]
	p.Home = parts[5]
	p.Shell = parts[6]

	return p, nil
}
