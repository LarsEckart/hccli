package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func AuthCmd() *cli.Command {
	return &cli.Command{
		Name:     "auth",
		Category: "Auth",
		Usage:    "Show API key info and permissions",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			auth, err := client.GetAuth(ctx)
			if err != nil {
				return err
			}

			return printJSON(auth)
		},
	}
}

func AuthV2Cmd() *cli.Command {
	return &cli.Command{
		Name:     "auth-v2",
		Category: "Auth",
		Usage:    "Show management API key info and permissions (v2)",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			auth, err := client.GetAuthV2(ctx)
			if err != nil {
				return err
			}

			return printJSON(auth)
		},
	}
}
