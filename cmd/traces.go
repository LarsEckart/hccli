package cmd

import (
	"context"
	"fmt"
	"net/url"

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
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			auth, err := client.GetAuth(ctx)
			if err != nil {
				return fmt.Errorf("fetching auth info: %w", err)
			}

			traceID := cmd.String("trace-id")
			dataset := cmd.String("dataset")
			teamSlug := auth.Team.Slug
			envSlug := auth.Environment.Slug

			result := map[string]string{
				"trace_id":    traceID,
				"dataset":     dataset,
				"team":        teamSlug,
				"environment": envSlug,
				"url":         traceUIURL(teamSlug, envSlug, dataset, traceID),
			}

			return printJSON(result)
		},
	}
}

func traceUIURL(teamSlug string, envSlug string, dataset string, traceID string) string {
	return fmt.Sprintf(
		"https://ui.honeycomb.io/%s/environments/%s/datasets/%s/trace?%s",
		url.PathEscape(teamSlug),
		url.PathEscape(envSlug),
		url.PathEscape(dataset),
		url.Values{"trace_id": []string{traceID}}.Encode(),
	)
}
