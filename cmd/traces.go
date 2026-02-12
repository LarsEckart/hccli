package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func GetTraceCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-trace",
		Category: "Traces",
		Usage:    "Get the Honeycomb UI URL for a trace",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "trace-id",
				Usage:    "Trace ID to look up",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			auth, err := client.GetAuth(ctx)
			if err != nil {
				return fmt.Errorf("fetching auth info: %w", err)
			}

			traceID := cmd.String("trace-id")
			dataset := cmd.String("dataset")
			teamSlug := auth.Team.Slug
			envSlug := auth.Environment.Slug

			traceURL := fmt.Sprintf(
				"https://ui.honeycomb.io/%s/environments/%s/datasets/%s/trace?trace_id=%s",
				teamSlug, envSlug, dataset, traceID,
			)

			result := map[string]string{
				"trace_id":    traceID,
				"dataset":     dataset,
				"team":        teamSlug,
				"environment": envSlug,
				"url":         traceURL,
			}

			return printJSON(result)
		},
	}
}
