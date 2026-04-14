package main

import (
	"fmt"
	"os"

	"github.com/aaangelmartin/goto/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "goto:", err)
		os.Exit(1)
	}
}
