package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

// noValueOps are filter operators that take no value argument.
var noValueOps = map[string]bool{
	"exists":         true,
	"does-not-exist": true,
}

// parseFilter parses a filter string in the form "column op [value]".
// The column is the first whitespace-delimited token, the op is the second,
// and the optional value is everything after the op.
func parseFilter(s string) (api.QueryFilter, error) {
	// Split into at most 3 parts: column, op, value
	parts := strings.SplitN(strings.TrimSpace(s), " ", 3)
	if len(parts) < 2 {
		return api.QueryFilter{}, fmt.Errorf("invalid filter %q: expected \"column op [value]\"", s)
	}

	col := parts[0]
	op := parts[1]

	f := api.QueryFilter{
		Column: col,
		Op:     op,
	}

	if noValueOps[op] {
		return f, nil
	}

	if len(parts) < 3 || parts[2] == "" {
		return api.QueryFilter{}, fmt.Errorf("invalid filter %q: operator %q requires a value", s, op)
	}
	f.Value = parts[2]
	return f, nil
}

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
			&cli.StringSliceFlag{
				Name:  "filter",
				Usage: `Filter in "column op [value]" form; repeat for multiple filters (e.g. --filter "duration_ms > 100" --filter "name exists")`,
			},
			&cli.StringFlag{
				Name:  "filter-combination",
				Usage: "How to combine filters: AND (default) or OR",
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

			for _, raw := range cmd.StringSlice("filter") {
				f, err := parseFilter(raw)
				if err != nil {
					return err
				}
				query.Filters = append(query.Filters, f)
			}

			if v := cmd.String("filter-combination"); v != "" {
				query.FilterCombination = v
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
