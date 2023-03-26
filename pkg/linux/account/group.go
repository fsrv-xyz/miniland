package account

import (
	"errors"
)

type GroupOption func(*Group) error

type Group struct {
	groupEntry
}

func GroupWithName(name string) GroupOption {
	return func(group *Group) error {
		// validate the name
		if !isValidLinuxGroupName(name) {
			return errors.New("invalid group name")
		}
		group.Name = name
		return nil
	}
}

func GroupWithGid(gid int) GroupOption {
	return func(group *Group) error {
		group.Gid = gid
		return nil
	}
}

func GroupWithMembers(members []string) GroupOption {
	return func(group *Group) error {
		group.Members = members
		return nil
	}
}
