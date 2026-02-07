package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "hccli",
		Usage: "Honeycomb API CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "api-key",
				Sources: cli.EnvVars("HONEYCOMB_API_KEY"),
				Usage:   "Honeycomb API key",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("hccli - Honeycomb API CLI")
			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
