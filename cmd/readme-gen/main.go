package main

import (
	"os"

	"github.com/hulk510/readme-gen/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
