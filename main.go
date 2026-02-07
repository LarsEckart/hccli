package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "hccli",
		Usage: "Honeycomb API CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "api-key",
				Sources:  cli.EnvVars("HONEYCOMB_API_KEY"),
				Usage:    "Honeycomb API key",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			authCmd(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func newClient(cmd *cli.Command) *api.Client {
	return api.NewClient(cmd.String("api-key"))
}

func authCmd() *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "Show API key info and permissions",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			auth, err := client.GetAuth(ctx)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(auth)
		},
	}
}
