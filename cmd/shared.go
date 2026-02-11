package cmd

import (
	"time"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func newClient(cmd *cli.Command) *api.Client {
	timeout := time.Duration(cmd.Int("timeout")) * time.Second
	client := api.NewClient(cmd.String("api-key"), timeout)
	if url := cmd.String("api-url"); url != "" {
		client.BaseURL = url
	}
	return client
}

// IDFlag returns a standard ID flag.
func IDFlag(name, usage string) cli.Flag {
	return &cli.StringFlag{
		Name:     name,
		Usage:    usage,
		Required: true,
	}
}

// DatasetFlag returns the standard dataset flag.
func DatasetFlag() cli.Flag {
	return &cli.StringFlag{
		Name:     "dataset",
		Usage:    "Dataset slug (use __all__ for environment-wide)",
		Required: true,
	}
}
