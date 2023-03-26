package account

import (
	"regexp"
)

const maxGroupNameLength = 16
const maxUserNameLength = 32

func isValidLinuxUserName(userName string) bool {
	// Check if the username is empty or too long
	if len(userName) == 0 || len(userName) > maxUserNameLength {
		return false
	}

	// Use regular expression to validate user name
	return regexp.MustCompile(`^[a-z][a-z0-9_-]*$`).MatchString(userName)
}

func isValidLinuxGroupName(groupName string) bool {
	// Check if the group name is empty or too long
	if groupName == "" || len(groupName) > maxGroupNameLength {
		return false
	}

	// Use regular expression to validate group name
	return regexp.MustCompile("^[a-z_][a-z0-9_-]*$").MatchString(groupName)
}
