package account

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePasswdLine(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		input    string
		expected any
		err      error
	}{
		{
			name:  "valid entry 1",
			input: "user1:x:1000:1000:User One:/home/user1:/bin/bash",
			expected: &passwdEntry{
				Username: "user1",
				Password: "x",
				Uid:      1000,
				Gid:      1000,
				Gecos:    "User One",
				Home:     "/home/user1",
				Shell:    "/bin/bash",
			},
			err: nil,
		},
		{
			name:  "valid entry 2",
			input: "user2:x:1001:1001:User Two:/home/user2:/bin/zsh",
			expected: &passwdEntry{
				Username: "user2",
				Password: "x",
				Uid:      1001,
				Gid:      1001,
				Gecos:    "User Two",
				Home:     "/home/user2",
				Shell:    "/bin/zsh",
			},
			err: nil,
		},
		{
			name:     "invalid entry",
			input:    "invalid:entry",
			expected: nil,
			err:      errors.New("invalid account entry"),
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			entry, err := passwdEntry{}.parse(testCase.input)
			assert.Equal(t, testCase.expected, entry)
			assert.Equal(t, testCase.err, err)
		})
	}
}

func TestPasswdEntry_String(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		entry    *passwdEntry
		expected string
	}{
		{
			name: "valid entry 1",
			entry: &passwdEntry{
				Username: "user1",
				Password: "x",
				Uid:      1000,
				Gid:      1000,
				Gecos:    "User One",
				Home:     "/home/user1",
				Shell:    "/bin/bash",
			},
			expected: "user1:x:1000:1000:User One:/home/user1:/bin/bash",
		},
		{
			name:     "empty entry",
			entry:    &passwdEntry{},
			expected: "",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.entry.String())
		})
	}
}
