package main

import (
	"log"
	"os"

	"github.com/alexalexdidit/secretmanager/internal/cli"
)

func main() {
	cmd := cli.NewRootCommand()

	if err := cmd.Execute(); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}
