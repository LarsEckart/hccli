package cmd

import (
	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func newClient(cmd *cli.Command) *api.Client {
	return api.NewClient(cmd.String("api-key"))
}
