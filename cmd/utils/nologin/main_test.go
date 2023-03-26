package main

import (
	"os"
	"testing"
)

func TestGetNoLoginMessage(t *testing.T) {
	expectedDefaultMessage := "This account is currently not available."

	message := getNoLoginMessage()
	if message != expectedDefaultMessage {
		t.Errorf("Expected default message '%s', got '%s'", expectedDefaultMessage, message)
	}

	customMessage := "Custom no login message."
	os.Setenv("GO_NOLOGIN_MESSAGE", customMessage)

	message = getNoLoginMessage()
	if message != customMessage {
		t.Errorf("Expected custom message '%s', got '%s'", customMessage, message)
	}

	os.Unsetenv("GO_NOLOGIN_MESSAGE")
}
