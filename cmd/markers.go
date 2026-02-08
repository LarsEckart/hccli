package cmd

import (
	"context"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func ListMarkersCmd() *cli.Command {
	return &cli.Command{
		Name:     "markers",
		Category: "Markers",
		Usage:    "List all markers for a dataset",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			markers, err := client.ListMarkers(ctx, cmd.String("dataset"))
			if err != nil {
				return err
			}

			return printJSON(markers)
		},
	}
}

func CreateMarkerCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-marker",
		Category: "Markers",
		Usage:    "Create a marker in a dataset",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "message",
				Usage: "Marker message",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "Marker type (e.g. deploy)",
			},
			&cli.StringFlag{
				Name:  "url",
				Usage: "URL target for the marker",
			},
			&cli.IntFlag{
				Name:  "start-time",
				Usage: "Start time as Unix timestamp",
			},
			&cli.IntFlag{
				Name:  "end-time",
				Usage: "End time as Unix timestamp",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			m := &api.Marker{
				Message: cmd.String("message"),
				Type:    cmd.String("type"),
				URL:     cmd.String("url"),
			}
			if cmd.IsSet("start-time") {
				v := int64(cmd.Int("start-time"))
				m.StartTime = &v
			}
			if cmd.IsSet("end-time") {
				v := int64(cmd.Int("end-time"))
				m.EndTime = &v
			}

			created, err := client.CreateMarker(ctx, cmd.String("dataset"), m)
			if err != nil {
				return err
			}

			return printJSON(created)
		},
	}
}

func UpdateMarkerCmd() *cli.Command {
	return &cli.Command{
		Name:     "update-marker",
		Category: "Markers",
		Usage:    "Update a marker by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Marker ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "message",
				Usage: "Marker message",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "Marker type (e.g. deploy)",
			},
			&cli.StringFlag{
				Name:  "url",
				Usage: "URL target for the marker",
			},
			&cli.IntFlag{
				Name:  "start-time",
				Usage: "Start time as Unix timestamp",
			},
			&cli.IntFlag{
				Name:  "end-time",
				Usage: "End time as Unix timestamp",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			m := &api.Marker{
				Message: cmd.String("message"),
				Type:    cmd.String("type"),
				URL:     cmd.String("url"),
			}
			if cmd.IsSet("start-time") {
				v := int64(cmd.Int("start-time"))
				m.StartTime = &v
			}
			if cmd.IsSet("end-time") {
				v := int64(cmd.Int("end-time"))
				m.EndTime = &v
			}

			updated, err := client.UpdateMarker(ctx, cmd.String("dataset"), cmd.String("id"), m)
			if err != nil {
				return err
			}

			return printJSON(updated)
		},
	}
}

func DeleteMarkerCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete-marker",
		Category: "Markers",
		Usage:    "Delete a marker by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Marker ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			return client.DeleteMarker(ctx, cmd.String("dataset"), cmd.String("id"))
		},
	}
}
