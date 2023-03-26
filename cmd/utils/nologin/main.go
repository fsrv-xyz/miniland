package main

import (
	"fmt"
	"os"
)

func getNoLoginMessage() string {
	_, exists := os.LookupEnv("GO_NOLOGIN_MESSAGE")
	message := "This account is currently not available."

	if exists {
		return os.Getenv("GO_NOLOGIN_MESSAGE")
	}

	return message
}

func main() {
	fmt.Println(getNoLoginMessage())
	os.Exit(1)
}
