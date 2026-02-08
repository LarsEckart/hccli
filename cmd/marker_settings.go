package cmd

import (
	"context"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func ListMarkerSettingsCmd() *cli.Command {
	return &cli.Command{
		Name:     "marker-settings",
		Category: "Marker Settings",
		Usage:    "List all marker settings for a dataset",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			settings, err := client.ListMarkerSettings(ctx, cmd.String("dataset"))
			if err != nil {
				return err
			}

			return printJSON(settings)
		},
	}
}

func CreateMarkerSettingCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-marker-setting",
		Category: "Marker Settings",
		Usage:    "Create a marker setting in a dataset",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "type",
				Usage:    "Marker type (e.g. deploy)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "color",
				Usage:    "Color as hexadecimal RGB (e.g. #FF0000)",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			ms := &api.MarkerSetting{
				Type:  cmd.String("type"),
				Color: cmd.String("color"),
			}

			created, err := client.CreateMarkerSetting(ctx, cmd.String("dataset"), ms)
			if err != nil {
				return err
			}

			return printJSON(created)
		},
	}
}

func UpdateMarkerSettingCmd() *cli.Command {
	return &cli.Command{
		Name:     "update-marker-setting",
		Category: "Marker Settings",
		Usage:    "Update a marker setting by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Marker setting ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "type",
				Usage:    "Marker type (e.g. deploy)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "color",
				Usage:    "Color as hexadecimal RGB (e.g. #FF0000)",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			ms := &api.MarkerSetting{
				Type:  cmd.String("type"),
				Color: cmd.String("color"),
			}

			updated, err := client.UpdateMarkerSetting(ctx, cmd.String("dataset"), cmd.String("id"), ms)
			if err != nil {
				return err
			}

			return printJSON(updated)
		},
	}
}

func DeleteMarkerSettingCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete-marker-setting",
		Category: "Marker Settings",
		Usage:    "Delete a marker setting by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Marker setting ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			return client.DeleteMarkerSetting(ctx, cmd.String("dataset"), cmd.String("id"))
		},
	}
}
