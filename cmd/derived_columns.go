package cmd

import (
	"context"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func ListDerivedColumnsCmd() *cli.Command {
	return &cli.Command{
		Name:     "derived-columns",
		Category: "Derived Columns",
		Usage:    "List all calculated fields (derived columns)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			cols, err := client.ListDerivedColumns(ctx, cmd.String("dataset"))
			if err != nil {
				return err
			}

			return printJSON(cols)
		},
	}
}

func GetDerivedColumnCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-derived-column",
		Category: "Derived Columns",
		Usage:    "Get a calculated field (derived column) by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Derived column ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			col, err := client.GetDerivedColumn(ctx, cmd.String("dataset"), cmd.String("id"))
			if err != nil {
				return err
			}

			return printJSON(col)
		},
	}
}

func CreateDerivedColumnCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-derived-column",
		Category: "Derived Columns",
		Usage:    "Create a new calculated field (derived column)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "alias",
				Usage:    "Human-readable name for the calculated field",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "expression",
				Usage:    "Expression to evaluate (e.g. INT(1))",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Human-readable description",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			col := &api.DerivedColumn{
				Alias:      cmd.String("alias"),
				Expression: cmd.String("expression"),
			}
			if v := cmd.String("description"); v != "" {
				col.Description = v
			}

			created, err := client.CreateDerivedColumn(ctx, cmd.String("dataset"), col)
			if err != nil {
				return err
			}

			return printJSON(created)
		},
	}
}

func UpdateDerivedColumnCmd() *cli.Command {
	return &cli.Command{
		Name:     "update-derived-column",
		Category: "Derived Columns",
		Usage:    "Update a calculated field (derived column) by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Derived column ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "alias",
				Usage:    "Human-readable name for the calculated field",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "expression",
				Usage:    "Expression to evaluate",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Human-readable description",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			col := &api.DerivedColumn{
				Alias:      cmd.String("alias"),
				Expression: cmd.String("expression"),
			}
			if v := cmd.String("description"); v != "" {
				col.Description = v
			}

			updated, err := client.UpdateDerivedColumn(ctx, cmd.String("dataset"), cmd.String("id"), col)
			if err != nil {
				return err
			}

			return printJSON(updated)
		},
	}
}

func DeleteDerivedColumnCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete-derived-column",
		Category: "Derived Columns",
		Usage:    "Delete a calculated field (derived column) by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Derived column ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			return client.DeleteDerivedColumn(ctx, cmd.String("dataset"), cmd.String("id"))
		},
	}
}
