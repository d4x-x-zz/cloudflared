package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

// TunnelCommand returns the CLI command for managing tunnels.
func TunnelCommand() *cli.Command {
	return &cli.Command{
		Name:    "tunnel",
		Aliases: []string{"t"},
		Usage:   "Use Cloudflare Tunnel to expose private services to the internet or to Cloudflare connected private networks.",
		Subcommands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Start a Cloudflare Tunnel",
				Action: runTunnel,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "Path to the configuration file",
						EnvVars: []string{"TUNNEL_CONFIG"},
					},
					&cli.StringFlag{
						Name:    "token",
						Usage:   "The Tunnel token",
						EnvVars: []string{"TUNNEL_TOKEN"},
					},
					&cli.BoolFlag{
						Name:  "no-autoupdate",
						Usage: "Disable automatic updates",
						// Default to true so autoupdates don't surprise me on my personal machines
						Value: true,
					},
				},
			},
			{
				Name:   "list",
				Usage:  "List existing tunnels",
				Action: listTunnels,
			},
			{
				Name:   "create",
				Usage:  "Create a new tunnel",
				Action: createTunnel,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "Name for the new tunnel",
						Required: true,
					},
				},
			},
			{
				Name:   "delete",
				Usage:  "Delete an existing tunnel",
				Action: deleteTunnel,
			},
		},
	}
}

// initLogger sets up the global logger with the given log level.
func initLogger(level string) zerolog.Logger {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		parsedLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(parsedLevel)

	return log.Logger
}

func runTunnel(c *cli.Context) error {
	logger := initLogger(c.String("loglevel"))

	token := c.String("token")
	if token == "" {
		return fmt.Errorf("tunnel token is required; set --token or TUNNEL_TOKEN env var")
	}

	logger.Info().Msg("Starting Cloudflare Tunnel...")
	logger.Info().Str("config", c.String("config")).Msg("Using configuration")

	// TODO: initialize and run the actual tunnel connector
	return fmt.Errorf("tunnel run not yet implemented")
}

func listTunnels(c *cli.Context) error {
	logger := initLogger(c.String("loglevel"))
	logger.Info().Msg("Listing tunnels...")

	// TODO: fetch and display tunnels from Cloudflare API
	fmt.Println("No tunnels found.")
	return nil
}

func createTunnel(c *cli.Context) error {
	logger := initLogger(c.String("loglevel"))
	name := c.String("name")

	logger.Info().Str("name", name).Msg("Creating tunnel...")

	// TODO: call Cloudflare API to create the tunnel
	return fmt.Errorf("tunnel create not yet implemented")
}

func deleteTunnel(c *cli.Context) error {
	logger := initLogger(c.String("loglevel"))

	if c.NArg() == 0 {
		return fmt.Errorf("tunnel ID or name is required")
	}
	tunnelID := c.Args().First()

	logger.Info().Str