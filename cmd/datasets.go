package cmd

import (
	"context"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func ListDatasetsCmd() *cli.Command {
	return &cli.Command{
		Name:     "datasets",
		Category: "Datasets",
		Usage:    "List all datasets",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			datasets, err := client.ListDatasets(ctx)
			if err != nil {
				return err
			}

			return printJSON(datasets)
		},
	}
}

func GetDatasetCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-dataset",
		Category: "Datasets",
		Usage:    "Get a dataset by slug",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "slug",
				Usage:    "Dataset slug",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			ds, err := client.GetDataset(ctx, cmd.String("slug"))
			if err != nil {
				return err
			}

			return printJSON(ds)
		},
	}
}

func CreateDatasetCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-dataset",
		Category: "Datasets",
		Usage:    "Create a new dataset",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Dataset name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Dataset description",
			},
			&cli.IntFlag{
				Name:  "expand-json-depth",
				Usage: "Maximum unpacking depth of nested JSON fields (0-10)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			ds := &api.Dataset{
				Name: cmd.String("name"),
			}
			if v := cmd.String("description"); v != "" {
				ds.Description = v
			}
			if cmd.IsSet("expand-json-depth") {
				depth := int(cmd.Int("expand-json-depth"))
				ds.ExpandJSONDepth = &depth
			}

			created, err := client.CreateDataset(ctx, ds)
			if err != nil {
				return err
			}

			return printJSON(created)
		},
	}
}

func UpdateDatasetCmd() *cli.Command {
	return &cli.Command{
		Name:     "update-dataset",
		Category: "Datasets",
		Usage:    "Update a dataset by slug",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "slug",
				Usage:    "Dataset slug",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "description",
				Usage:    "Dataset description",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "expand-json-depth",
				Usage:    "Maximum unpacking depth of nested JSON fields (0-10)",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "delete-protected",
				Usage: "Enable deletion protection",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			depth := int(cmd.Int("expand-json-depth"))
			ds := &api.Dataset{
				Description:     cmd.String("description"),
				ExpandJSONDepth: &depth,
			}
			if cmd.IsSet("delete-protected") {
				dp := cmd.Bool("delete-protected")
				ds.Settings = &api.DatasetSettings{
					DeleteProtected: &dp,
				}
			}

			updated, err := client.UpdateDataset(ctx, cmd.String("slug"), ds)
			if err != nil {
				return err
			}

			return printJSON(updated)
		},
	}
}

func DeleteDatasetCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete-dataset",
		Category: "Datasets",
		Usage:    "Delete a dataset by slug",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "slug",
				Usage:    "Dataset slug",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			return client.DeleteDataset(ctx, cmd.String("slug"))
		},
	}
}
