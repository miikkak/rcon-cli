// Command rcon-cli is an interactive and single-shot CLI for RCON-enabled
// game servers, such as Minecraft.
package main

import (
	"fmt"
	"os"

	"github.com/miikkak/rcon-cli/internal/cmd"
)

var version = "dev"

func main() {
	cmd.RootCmd.Version = version
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
