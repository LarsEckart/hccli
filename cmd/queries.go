package cmd

import (
	"context"
	"fmt"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func GetQueryCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-query",
		Category: "Queries",
		Usage:    "Get a query by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Query ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			query, err := client.GetQuery(ctx, cmd.String("dataset"), cmd.String("id"))
			if err != nil {
				return err
			}

			return printJSON(query)
		},
	}
}

func CreateQueryCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-query",
		Category: "Queries",
		Usage:    "Create a new query",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "calculation-op",
				Usage:    "Calculation operation (e.g. COUNT, AVG, P99); repeat for multiple calculations",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:  "calculation-column",
				Usage: "Calculation column (use empty string for COUNT); repeat to match each --calculation-op",
			},
			&cli.StringFlag{
				Name:  "breakdown",
				Usage: "Breakdown column",
			},
			&cli.StringFlag{
				Name:  "filter-column",
				Usage: "Filter column",
			},
			&cli.StringFlag{
				Name:  "filter-op",
				Usage: "Filter operation",
			},
			&cli.StringFlag{
				Name:  "filter-value",
				Usage: "Filter value",
			},
			&cli.IntFlag{
				Name:  "time-range",
				Usage: "Time range in seconds",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			ops := cmd.StringSlice("calculation-op")
			cols := cmd.StringSlice("calculation-column")

			if len(cols) > 0 && len(cols) != len(ops) {
				return fmt.Errorf("number of --calculation-column values (%d) must match --calculation-op values (%d)", len(cols), len(ops))
			}

			var calcs []api.Calculation
			for i, op := range ops {
				c := api.Calculation{Op: op}
				if i < len(cols) && cols[i] != "" {
					c.Column = cols[i]
				}
				calcs = append(calcs, c)
			}

			query := &api.Query{
				Calculations: calcs,
			}

			if v := cmd.String("breakdown"); v != "" {
				query.Breakdowns = []string{v}
			}

			if col := cmd.String("filter-column"); col != "" {
				f := api.QueryFilter{
					Column: col,
					Op:     cmd.String("filter-op"),
				}
				if v := cmd.String("filter-value"); v != "" {
					f.Value = v
				}
				query.Filters = []api.QueryFilter{f}
			}

			if v := cmd.Int("time-range"); v != 0 {
				query.TimeRange = int(v)
			}

			created, err := client.CreateQuery(ctx, cmd.String("dataset"), query)
			if err != nil {
				return err
			}

			return printJSON(created)
		},
	}
}
