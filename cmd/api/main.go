package main

import (
	"log"

	"github.com/copilot-ghostpilot-app/api"
)

func main() {
	if err := api.Run(); err != nil {
		log.Fatalf("run api server: %v\n", err)
	}
}
