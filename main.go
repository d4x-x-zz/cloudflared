package main

import (
	"fmt"
	"os"

	"github.com/cloudflare/cloudflared/cmd/cloudflared"
)

// main is the entry point for cloudflared.
// It delegates to the cloudflared command package which sets up
// the CLI application and handles subcommands.
//
// Personal fork - using this for learning how Cloudflare tunnels work
// and experimenting with custom tunnel configurations.
func main() {
	if err := cloudflared.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
