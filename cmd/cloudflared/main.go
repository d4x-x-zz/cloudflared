package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/cloudflare/cloudflared/cmd/cloudflared/access"
	"github.com/cloudflare/cloudflared/cmd/cloudflared/tunnel"
	"github.com/cloudflare/cloudflared/cmd/cloudflared/updater"
)

var (
	// Version is the version of cloudflared, injected at build time
	Version = "DEV"
	// BuildTime is the time cloudflared was built, injected at build time
	BuildTime = "unknown"
)

func main() {
	app := &cli.App{
		Name:    "cloudflared",
		Usage:   "Cloudflare Tunnel client",
		Version: fmt.Sprintf("%s (built %s)", Version, BuildTime),
		Authors: []*cli.Author{
			{
				Name:  "Cloudflare",
				Email: "support@cloudflare.com",
			},
		},
		Commands: []*cli.Command{
			tunnel.Commands(),
			access.Commands(),
			updater.Commands(),
		},
		// Default action runs tunnel when no subcommand is given
		Action: tunnel.TunnelCommand,
		Flags:  tunnel.Flags(),
		Before: func(c *cli.Context) error {
			return nil
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				// print to stderr and exit with code 1 so scripts can detect failures easily
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
		// suggest subcommands if the user makes a typo
		SuggestAfterError: true,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
