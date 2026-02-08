package cmd

import (
	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func newClient(cmd *cli.Command) *api.Client {
	return api.NewClient(cmd.String("api-key"))
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
