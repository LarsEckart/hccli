package cmd

import (
	"context"
	"time"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

type deleteOutput struct {
	Deleted  bool   `json:"deleted"`
	Resource string `json:"resource"`
	ID       string `json:"id,omitzero"`
	Slug     string `json:"slug,omitzero"`
	Dataset  string `json:"dataset,omitzero"`
	BoardID  string `json:"board_id,omitzero"`
	ViewID   string `json:"view_id,omitzero"`
}

func printDeleteResult(resource string, fields deleteOutput) error {
	fields.Deleted = true
	fields.Resource = resource
	return printJSON(fields)
}

func newClient(cmd *cli.Command) (*api.Client, error) {
	credentials, err := resolveCredentials(cmd)
	if err != nil {
		return nil, err
	}
	timeout := time.Duration(cmd.Root().Int("timeout")) * time.Second
	client := api.NewClient(credentials.APIKey, timeout)
	if credentials.APIURL != "" {
		client.BaseURL = credentials.APIURL
	}
	return client, nil
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

func deleteAction(resource string, deleteFn func(context.Context, *cli.Command) error, fieldsFn func(*cli.Command) deleteOutput) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		if err := deleteFn(ctx, cmd); err != nil {
			return err
		}
		return printDeleteResult(resource, fieldsFn(cmd))
	}
}
