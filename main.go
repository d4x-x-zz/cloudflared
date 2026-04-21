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
//
// Note: if the process exits with a non-zero code, the error is printed
// to stderr before exiting so it's easier to spot in logs/terminal output.
// Also printing a newline before the error so it stands out more when
// there's a bunch of log output scrolling by.
//
// TODO: look into how the reconnect backoff logic works in the tunnel package,
// might want to tweak the retry intervals for my home server setup.
func main() {
	if err := cloudflared.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\nerror: %v\n", err)
		os.Exit(1)
	}
}
