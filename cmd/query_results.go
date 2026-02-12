package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func CreateQueryResultCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-query-result",
		Category: "Query Results",
		Usage:    "Execute a query and return results (polls until complete)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "query-id",
				Usage:    "Query ID to execute",
				Required: true,
			},
			&cli.IntFlag{
				Name:  "poll-interval",
				Usage: "Seconds between polling attempts",
				Value: 2,
			},
			&cli.IntFlag{
				Name:  "timeout",
				Usage: "Maximum seconds to wait for results",
				Value: 60,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			dataset := cmd.String("dataset")
			queryID := cmd.String("query-id")
			pollInterval := time.Duration(cmd.Int("poll-interval")) * time.Second
			timeout := time.Duration(cmd.Int("timeout")) * time.Second

			if pollInterval < 1*time.Second {
				pollInterval = 1 * time.Second
			}

			result, err := client.CreateQueryResult(ctx, dataset, queryID)
			if err != nil {
				return err
			}

			if result.Complete {
				warnIfEmptyResults(result, dataset)
				return printJSON(result)
			}

			deadline := time.Now().Add(timeout)
			for !result.Complete {
				if time.Now().After(deadline) {
					return fmt.Errorf("timed out waiting for query result %s after %s", result.ID, timeout)
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(pollInterval):
				}

				result, err = client.GetQueryResult(ctx, dataset, result.ID)
				if err != nil {
					return err
				}
			}

			warnIfEmptyResults(result, dataset)
			return printJSON(result)
		},
	}
}

func warnIfEmptyResults(result *api.QueryResult, dataset string) {
	if len(result.Data.Results) > 0 {
		return
	}
	fmt.Fprintln(os.Stderr, "⚠️  Query returned 0 results")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "💡 Possible reasons:")
	fmt.Fprintln(os.Stderr, "  • No data in time range (try a larger --time-range)")
	fmt.Fprintln(os.Stderr, "  • Filters are too restrictive")
	fmt.Fprintf(os.Stderr, "  • Column names don't exist (verify with: hccli columns --dataset %s)\n", dataset)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "💡 Results with breakdowns are nested under .data.results[].data")
	fmt.Fprintln(os.Stderr, "   Try: jq '.data.results[].data'")
}

func GetQueryResultCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-query-result",
		Category: "Query Results",
		Usage:    "Get a query result by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Query result ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			result, err := client.GetQueryResult(ctx, cmd.String("dataset"), cmd.String("id"))
			if err != nil {
				return err
			}

			return printJSON(result)
		},
	}
}
