package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LarsEckart/hccli/cmd"
	"github.com/urfave/cli/v3"
)

func newApp() *cli.Command {
	cli.VersionFlag = &cli.BoolFlag{
		Name:        "version",
		Usage:       "print the version",
		HideDefault: true,
		Local:       true,
	}

	return &cli.Command{
		Name:    "hccli",
		Usage:   "A machine-friendly CLI for the Honeycomb observability platform",
		Version: appVersion(),
		Description: `Interact with Honeycomb from the command line — ideal for scripting,
automation, and integration with CI/CD pipelines.

Authentication:
  Run hccli auth login --profile <name> --api-key-stdin to store named profiles.
  Use hccli auth switch <name> to change the active profile.
  For CI and one-off use, pass --api-key or set HONEYCOMB_API_KEY.

Output:
  All commands output JSON with 2-space indentation, making them easy to parse
  and pipe into tools like jq for further processing.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "api-key",
				Sources: cli.EnvVars("HONEYCOMB_API_KEY"),
				Usage:   "Honeycomb API key (overrides profiles)",
			},
			&cli.StringFlag{
				Name:    "profile",
				Aliases: []string{"p"},
				Sources: cli.EnvVars("HCCLI_PROFILE"),
				Usage:   "Honeycomb auth profile to use",
			},
			&cli.IntFlag{
				Name:        "timeout",
				Usage:       "HTTP request timeout in seconds",
				Value:       30,
				DefaultText: "30",
				Sources:     cli.EnvVars("HONEYCOMB_TIMEOUT"),
			},
			&cli.StringFlag{
				Name:    "api-url",
				Usage:   "Honeycomb API base URL",
				Value:   "https://api.honeycomb.io",
				Sources: cli.EnvVars("HONEYCOMB_API_URL"),
			},
		},
		Commands: []*cli.Command{
			cmd.AuthCmd(),
			cmd.AuthV2Cmd(),
			cmd.ListBoardsCmd(),
			cmd.GetBoardCmd(),
			cmd.CreateBoardCmd(),
			cmd.UpdateBoardCmd(),
			cmd.DeleteBoardCmd(),
			cmd.ListBoardViewsCmd(),
			cmd.GetBoardViewCmd(),
			cmd.CreateBoardViewCmd(),
			cmd.UpdateBoardViewCmd(),
			cmd.DeleteBoardViewCmd(),
			cmd.GetQueryCmd(),
			cmd.CreateQueryCmd(),
			cmd.CreateQueryResultCmd(),
			cmd.GetQueryResultCmd(),
			cmd.CreateQueryAnnotationCmd(),
			cmd.ListQueryAnnotationsCmd(),
			cmd.GetQueryAnnotationCmd(),
			cmd.UpdateQueryAnnotationCmd(),
			cmd.DeleteQueryAnnotationCmd(),
			cmd.ListColumnsCmd(),
			cmd.GetColumnCmd(),
			cmd.CreateColumnCmd(),
			cmd.UpdateColumnCmd(),
			cmd.DeleteColumnCmd(),
			cmd.ListDatasetsCmd(),
			cmd.GetDatasetCmd(),
			cmd.CreateDatasetCmd(),
			cmd.UpdateDatasetCmd(),
			cmd.DeleteDatasetCmd(),
			cmd.ListDerivedColumnsCmd(),
			cmd.GetDerivedColumnCmd(),
			cmd.CreateDerivedColumnCmd(),
			cmd.UpdateDerivedColumnCmd(),
			cmd.DeleteDerivedColumnCmd(),
			cmd.ListMarkersCmd(),
			cmd.CreateMarkerCmd(),
			cmd.UpdateMarkerCmd(),
			cmd.DeleteMarkerCmd(),
			cmd.ListMarkerSettingsCmd(),
			cmd.CreateMarkerSettingCmd(),
			cmd.UpdateMarkerSettingCmd(),
			cmd.DeleteMarkerSettingCmd(),
			cmd.ListSLOsCmd(),
			cmd.GetSLOCmd(),
			cmd.CreateSLOCmd(),
			cmd.UpdateSLOCmd(),
			cmd.DeleteSLOCmd(),
			cmd.ListBurnAlertsCmd(),
			cmd.GetBurnAlertCmd(),
			cmd.CreateBurnAlertCmd(),
			cmd.UpdateBurnAlertCmd(),
			cmd.DeleteBurnAlertCmd(),
			cmd.GetTraceCmd(),
		},
	}
}

func main() {
	if err := newApp().Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
