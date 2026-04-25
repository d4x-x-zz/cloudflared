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
//
// TODO: also worth investigating whether I can reduce the number of connections
// (currently defaults to 4) since my home server is pretty low-traffic and
// 4 feels like overkill - probably 2 would be fine.
func main() {
	if err := cloudflared.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\nerror: %v\n", err)
		// print a separator so the error is easier to spot when
		// there's a wall of log output above it
		fmt.Fprintln(os.Stderr, "---")
		os.Exit(1)
	}
}
