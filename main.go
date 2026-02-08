package main

import (
	"context"
	"fmt"
	"os"

	"github.com/LarsEckart/hccli/cmd"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "hccli",
		Usage: "A machine-friendly CLI for the Honeycomb observability platform",
		Description: `Interact with Honeycomb from the command line — ideal for scripting,
automation, and integration with CI/CD pipelines.

Authentication:
  Provide your API key via --api-key flag or HONEYCOMB_API_KEY environment variable.

Output:
  All commands output JSON with 2-space indentation, making them easy to parse
  and pipe into tools like jq for further processing.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "api-key",
				Sources:  cli.EnvVars("HONEYCOMB_API_KEY"),
				Usage:    "Honeycomb API key",
				Required: true,
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
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
