package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidLinuxUserName(t *testing.T) {
	for _, testCase := range []struct {
		input    string
		expected bool
	}{
		{"valid_user", true},
		{"invalidUser", false},
		{"another_valid-user", true},
		{"invalid!user", false},
		{"1nvalid_user", false},
		{"a_very_long_username_that_exceeds_thirty_two_characters", false},
		{"", false},
	} {
		result := isValidLinuxUserName(testCase.input)
		assert.Equal(t, testCase.expected, result, "Expected isValidLinuxUserName(%q) to be %t, but got %t", testCase.input, testCase.expected, result)
	}
}

func TestIsValidLinuxGroupName(t *testing.T) {
	for _, testCase := range []struct {
		input    string
		expected bool
	}{
		{"valid_group", true},
		{"invalidGroup", false},
		{"too-long_group-name", false},
		{"invalid!group", false},
		{"1nvalid_group", false},
		{"a_very_long_groupname_that_exceeds_thirty_two_characters", false},
		{"", false},
	} {
		result := isValidLinuxGroupName(testCase.input)
		assert.Equal(t, testCase.expected, result, "Expected isValidLinuxGroupName(%q) to be %t, but got %t", testCase.input, testCase.expected, result)
	}
}
