package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func ListSLOsCmd() *cli.Command {
	return &cli.Command{
		Name:     "slos",
		Category: "SLOs",
		Usage:    "List all SLOs for a dataset",
		Flags: []cli.Flag{
			DatasetFlag(),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			slos, err := client.ListSLOs(ctx, cmd.String("dataset"))
			if err != nil {
				return err
			}
			return printJSON(slos)
		},
	}
}

func GetSLOCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-slo",
		Category: "SLOs",
		Usage:    "Get an SLO by ID",
		Flags: []cli.Flag{
			DatasetFlag(),
			IDFlag("id", "SLO ID"),
			&cli.BoolFlag{
				Name:  "detailed",
				Usage: "Include reporting data (compliance, budget_remaining, status, burn_rate)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			var (
				slo *api.SLO
				err error
			)
			if cmd.Bool("detailed") {
				slo, err = client.GetSLODetailed(ctx, cmd.String("dataset"), cmd.String("id"))
			} else {
				slo, err = client.GetSLO(ctx, cmd.String("dataset"), cmd.String("id"))
			}
			if err != nil {
				return err
			}
			return printJSON(slo)
		},
	}
}

func CreateSLOCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-slo",
		Category: "SLOs",
		Usage:    "Create an SLO",
		Flags: []cli.Flag{
			DatasetFlag(),
			&cli.StringFlag{
				Name:     "name",
				Usage:    "SLO name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "SLO description",
			},
			&cli.StringFlag{
				Name:     "sli-alias",
				Usage:    "Alias of the derived column to use as the SLI",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "time-period-days",
				Usage:    "Time period in days over which the SLO is evaluated",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "target-per-million",
				Usage:    "Target success rate per million (e.g. 999000 = 99.9%)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "tags-json",
				Usage: `JSON array of tags, e.g. '[{"key":"team","value":"blue"}]'`,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			slo := &api.SLO{
				Name:             cmd.String("name"),
				Description:      cmd.String("description"),
				SLI:              api.SLOSLI{Alias: cmd.String("sli-alias")},
				TimePeriodDays:   int(cmd.Int("time-period-days")),
				TargetPerMillion: int(cmd.Int("target-per-million")),
			}

			if tj := cmd.String("tags-json"); tj != "" {
				var tags []api.Tag
				if err := json.Unmarshal([]byte(tj), &tags); err != nil {
					return fmt.Errorf("parsing tags-json: %w", err)
				}
				slo.Tags = tags
			}

			created, err := client.CreateSLO(ctx, cmd.String("dataset"), slo)
			if err != nil {
				return err
			}
			return printJSON(created)
		},
	}
}

func UpdateSLOCmd() *cli.Command {
	return &cli.Command{
		Name:     "update-slo",
		Category: "SLOs",
		Usage:    "Update an SLO by ID",
		Flags: []cli.Flag{
			DatasetFlag(),
			IDFlag("id", "SLO ID"),
			&cli.StringFlag{
				Name:     "name",
				Usage:    "SLO name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "SLO description",
			},
			&cli.StringFlag{
				Name:     "sli-alias",
				Usage:    "Alias of the derived column to use as the SLI",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "time-period-days",
				Usage:    "Time period in days over which the SLO is evaluated",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "target-per-million",
				Usage:    "Target success rate per million (e.g. 999000 = 99.9%)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "tags-json",
				Usage: `JSON array of tags, e.g. '[{"key":"team","value":"blue"}]'`,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			slo := &api.SLO{
				Name:             cmd.String("name"),
				Description:      cmd.String("description"),
				SLI:              api.SLOSLI{Alias: cmd.String("sli-alias")},
				TimePeriodDays:   int(cmd.Int("time-period-days")),
				TargetPerMillion: int(cmd.Int("target-per-million")),
			}

			if tj := cmd.String("tags-json"); tj != "" {
				var tags []api.Tag
				if err := json.Unmarshal([]byte(tj), &tags); err != nil {
					return fmt.Errorf("parsing tags-json: %w", err)
				}
				slo.Tags = tags
			}

			updated, err := client.UpdateSLO(ctx, cmd.String("dataset"), cmd.String("id"), slo)
			if err != nil {
				return err
			}
			return printJSON(updated)
		},
	}
}

func DeleteSLOCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete-slo",
		Category: "SLOs",
		Usage:    "Delete an SLO by ID",
		Flags: []cli.Flag{
			DatasetFlag(),
			IDFlag("id", "SLO ID"),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			return client.DeleteSLO(ctx, cmd.String("dataset"), cmd.String("id"))
		},
	}
}
