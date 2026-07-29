package main

import (
	"fmt"
	"os"

	"github.com/Fauxmen4/deps-check/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
