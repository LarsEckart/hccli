package cmd

import (
	"context"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func ListColumnsCmd() *cli.Command {
	return &cli.Command{
		Name:     "columns",
		Category: "Columns",
		Usage:    "List all columns",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			cols, err := client.ListColumns(ctx, cmd.String("dataset"))
			if err != nil {
				return err
			}

			return printJSON(cols)
		},
	}
}

func GetColumnCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-column",
		Category: "Columns",
		Usage:    "Get a column by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Column ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			col, err := client.GetColumn(ctx, cmd.String("dataset"), cmd.String("id"))
			if err != nil {
				return err
			}

			return printJSON(col)
		},
	}
}

func CreateColumnCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-column",
		Category: "Columns",
		Usage:    "Create a new column",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "key-name",
				Usage:    "Column name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "Column type (string, float, integer, boolean)",
				Value: "string",
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Column description",
			},
			&cli.BoolFlag{
				Name:  "hidden",
				Usage: "Hide column from autocomplete",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			hidden := cmd.Bool("hidden")
			col := &api.Column{
				KeyName: cmd.String("key-name"),
				Type:    cmd.String("type"),
				Hidden:  &hidden,
			}
			if v := cmd.String("description"); v != "" {
				col.Description = v
			}

			created, err := client.CreateColumn(ctx, cmd.String("dataset"), col)
			if err != nil {
				return err
			}

			return printJSON(created)
		},
	}
}

func UpdateColumnCmd() *cli.Command {
	return &cli.Command{
		Name:     "update-column",
		Category: "Columns",
		Usage:    "Update a column by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Column ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "Column type (string, float, integer, boolean)",
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Column description",
			},
			&cli.BoolFlag{
				Name:  "hidden",
				Usage: "Hide column from autocomplete",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			col := &api.Column{}
			if v := cmd.String("type"); v != "" {
				col.Type = v
			}
			if v := cmd.String("description"); v != "" {
				col.Description = v
			}
			if cmd.IsSet("hidden") {
				hidden := cmd.Bool("hidden")
				col.Hidden = &hidden
			}

			updated, err := client.UpdateColumn(ctx, cmd.String("dataset"), cmd.String("id"), col)
			if err != nil {
				return err
			}

			return printJSON(updated)
		},
	}
}

func DeleteColumnCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete-column",
		Category: "Columns",
		Usage:    "Delete a column by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Column ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			return client.DeleteColumn(ctx, cmd.String("dataset"), cmd.String("id"))
		},
	}
}
