package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupEntryString(t *testing.T) {
	entry := &groupEntry{
		Name:     "testgroup",
		Password: "x",
		Gid:      1000,
		Members:  []string{"user1", "user2"},
	}

	expected := "testgroup:x:1000:user1,user2"
	result := entry.String()
	assert.Equal(t, expected, result)
}

func TestGroupEntryStringEmptyMembers(t *testing.T) {
	entry := &groupEntry{
		Name:     "testgroup",
		Password: "x",
		Gid:      1000,
		Members:  []string{},
	}

	expected := "testgroup:x:1000:"
	result := entry.String()
	assert.Equal(t, expected, result)
}

func TestGroupEntryStringEmpty(t *testing.T) {
	entry := &groupEntry{}

	expected := ""
	result := entry.String()
	assert.Equal(t, expected, result)
}

func TestParseGroupsLineValid(t *testing.T) {
	line := "testgroup:x:1000:user1,user2"
	expected := &groupEntry{
		Name:     "testgroup",
		Password: "x",
		Gid:      1000,
		Members:  []string{"user1", "user2"},
	}

	result, err := groupEntry{}.parse(line)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestParseGroupsLineInvalid(t *testing.T) {
	testCases := []struct {
		name  string
		line  string
		error string
	}{
		{
			name:  "invalid line",
			line:  "testgroup:x:1000",
			error: "invalid line: testgroup:x:1000",
		},
		{
			name:  "invalid gid",
			line:  "testgroup:x:notanumber:user1,user2",
			error: "invalid line: testgroup:x:notanumber:user1,user2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := groupEntry{}.parse(tc.line)
			assert.Error(t, err)
			assert.Equal(t, tc.error, err.Error())
		})
	}
}
