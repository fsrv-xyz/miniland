package main

import (
	"os"

	"ref.ci/fsrvcorp/miniland/userland/pkg/web"
)

func main() {
	address := os.Getenv("WEB_ADDRESS")
	if address == "" {
		address = ":8080"
	}
	web.Start(address)
}
