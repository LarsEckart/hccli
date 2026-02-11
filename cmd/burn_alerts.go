package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func ListBurnAlertsCmd() *cli.Command {
	return &cli.Command{
		Name:     "burn-alerts",
		Category: "Burn Alerts",
		Usage:    "List all burn alerts for an SLO",
		Flags: []cli.Flag{
			DatasetFlag(),
			&cli.StringFlag{
				Name:     "slo-id",
				Usage:    "SLO ID to list burn alerts for",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			alerts, err := client.ListBurnAlerts(ctx, cmd.String("dataset"), cmd.String("slo-id"))
			if err != nil {
				return err
			}
			return printJSON(alerts)
		},
	}
}

func GetBurnAlertCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-burn-alert",
		Category: "Burn Alerts",
		Usage:    "Get a burn alert by ID",
		Flags: []cli.Flag{
			DatasetFlag(),
			IDFlag("id", "Burn Alert ID"),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			alert, err := client.GetBurnAlert(ctx, cmd.String("dataset"), cmd.String("id"))
			if err != nil {
				return err
			}
			return printJSON(alert)
		},
	}
}

func CreateBurnAlertCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-burn-alert",
		Category: "Burn Alerts",
		Usage:    "Create a burn alert for an SLO",
		Description: `Create a burn alert. Two alert types are supported:

  exhaustion_time (default): fires when SLO budget will be exhausted
    within the specified number of minutes.

  budget_rate: fires when budget drops by a threshold percentage
    within a time window.

Examples:

  # Exhaustion time alert
  hccli create-burn-alert --dataset mydata --slo-id abc123 \
    --exhaustion-minutes 120 \
    --recipients-json '[{"type":"email","target":"alerts@example.com"}]'

  # Budget rate alert
  hccli create-burn-alert --dataset mydata --slo-id abc123 \
    --alert-type budget_rate \
    --budget-rate-window-minutes 60 \
    --budget-rate-decrease-per-million 10000 \
    --recipients-json '[{"type":"email","target":"alerts@example.com"}]'`,
		Flags: []cli.Flag{
			DatasetFlag(),
			&cli.StringFlag{
				Name:     "slo-id",
				Usage:    "SLO ID to attach the burn alert to",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "alert-type",
				Usage: "Alert type: exhaustion_time or budget_rate",
				Value: "exhaustion_time",
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Description of the burn alert",
			},
			&cli.IntFlag{
				Name:  "exhaustion-minutes",
				Usage: "Minutes until budget exhaustion triggers alert (for exhaustion_time)",
			},
			&cli.IntFlag{
				Name:  "budget-rate-window-minutes",
				Usage: "Time window in minutes for budget rate calculation (for budget_rate, min 60)",
			},
			&cli.IntFlag{
				Name:  "budget-rate-decrease-per-million",
				Usage: "Budget decrease threshold per million (for budget_rate, 1-1000000; 10000 = 1%)",
			},
			&cli.StringFlag{
				Name:     "recipients-json",
				Usage:    `JSON array of recipients, e.g. '[{"id":"abc123"}]' or '[{"type":"email","target":"a@b.com"}]'`,
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			ba, err := buildBurnAlert(cmd)
			if err != nil {
				return err
			}
			ba.SLO = &api.BurnAlertSLO{ID: cmd.String("slo-id")}

			created, err := client.CreateBurnAlert(ctx, cmd.String("dataset"), ba)
			if err != nil {
				return err
			}
			return printJSON(created)
		},
	}
}

func UpdateBurnAlertCmd() *cli.Command {
	return &cli.Command{
		Name:     "update-burn-alert",
		Category: "Burn Alerts",
		Usage:    "Update a burn alert by ID",
		Flags: []cli.Flag{
			DatasetFlag(),
			IDFlag("id", "Burn Alert ID"),
			&cli.StringFlag{
				Name:  "alert-type",
				Usage: "Alert type: exhaustion_time or budget_rate",
				Value: "exhaustion_time",
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Description of the burn alert",
			},
			&cli.IntFlag{
				Name:  "exhaustion-minutes",
				Usage: "Minutes until budget exhaustion triggers alert (for exhaustion_time)",
			},
			&cli.IntFlag{
				Name:  "budget-rate-window-minutes",
				Usage: "Time window in minutes for budget rate calculation (for budget_rate, min 60)",
			},
			&cli.IntFlag{
				Name:  "budget-rate-decrease-per-million",
				Usage: "Budget decrease threshold per million (for budget_rate, 1-1000000; 10000 = 1%)",
			},
			&cli.StringFlag{
				Name:     "recipients-json",
				Usage:    `JSON array of recipients, e.g. '[{"id":"abc123"}]' or '[{"type":"email","target":"a@b.com"}]'`,
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			ba, err := buildBurnAlert(cmd)
			if err != nil {
				return err
			}

			updated, err := client.UpdateBurnAlert(ctx, cmd.String("dataset"), cmd.String("id"), ba)
			if err != nil {
				return err
			}
			return printJSON(updated)
		},
	}
}

func DeleteBurnAlertCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete-burn-alert",
		Category: "Burn Alerts",
		Usage:    "Delete a burn alert by ID",
		Flags: []cli.Flag{
			DatasetFlag(),
			IDFlag("id", "Burn Alert ID"),
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			return client.DeleteBurnAlert(ctx, cmd.String("dataset"), cmd.String("id"))
		},
	}
}

func buildBurnAlert(cmd *cli.Command) (*api.BurnAlert, error) {
	var recipients []api.NotificationRecipient
	if err := json.Unmarshal([]byte(cmd.String("recipients-json")), &recipients); err != nil {
		return nil, fmt.Errorf("parsing recipients-json: %w", err)
	}

	ba := &api.BurnAlert{
		AlertType:   cmd.String("alert-type"),
		Description: cmd.String("description"),
		Recipients:  recipients,
	}

	switch ba.AlertType {
	case "exhaustion_time":
		if cmd.IsSet("exhaustion-minutes") {
			v := int(cmd.Int("exhaustion-minutes"))
			ba.ExhaustionMinutes = &v
		}
	case "budget_rate":
		if cmd.IsSet("budget-rate-window-minutes") {
			v := int(cmd.Int("budget-rate-window-minutes"))
			ba.BudgetRateWindowMinutes = &v
		}
		if cmd.IsSet("budget-rate-decrease-per-million") {
			v := int(cmd.Int("budget-rate-decrease-per-million"))
			ba.BudgetRateDecreaseThresholdPerMillion = &v
		}
	}

	return ba, nil
}
