package account

import (
	"errors"
)

type UserOption func(*User) error

type User struct {
	createGroup  bool
	groupOptions []GroupOption
	passwdEntry
}

func UserWithName(name string) UserOption {
	return func(user *User) error {
		// validate the name
		if !isValidLinuxUserName(name) {
			return errors.New("invalid user name")
		}
		user.Username = name
		return nil
	}
}

func UserWithPassword(password string) UserOption {
	return func(user *User) error {
		user.Password = password
		return nil
	}
}

func UserWithUid(uid int) UserOption {
	return func(user *User) error {
		user.Uid = uid
		return nil
	}
}

func UserWithGid(gid int) UserOption {
	return func(user *User) error {
		user.Gid = gid
		return nil
	}
}

func UserWithComment(gecos string) UserOption {
	return func(user *User) error {
		user.Gecos = gecos
		return nil
	}
}

func UserWithHome(home string) UserOption {
	return func(user *User) error {
		user.Home = home
		return nil
	}
}

func UserWithShell(shell string) UserOption {
	return func(user *User) error {
		user.Shell = shell
		return nil
	}
}

func UserCreateGroup(groupOptions ...GroupOption) UserOption {
	return func(user *User) error {
		user.createGroup = true
		user.groupOptions = groupOptions
		return nil
	}
}
