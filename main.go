package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "hccli",
		Usage: "Honeycomb API CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "api-key",
				Sources:  cli.EnvVars("HONEYCOMB_API_KEY"),
				Usage:    "Honeycomb API key",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			authCmd(),
			authV2Cmd(),
			listBoardsCmd(),
			getBoardCmd(),
			createBoardCmd(),
			updateBoardCmd(),
			deleteBoardCmd(),
			listBoardViewsCmd(),
			getBoardViewCmd(),
			createBoardViewCmd(),
			updateBoardViewCmd(),
			deleteBoardViewCmd(),
			getQueryCmd(),
			createQueryCmd(),
			createQueryAnnotationCmd(),
			listDerivedColumnsCmd(),
			getDerivedColumnCmd(),
			createDerivedColumnCmd(),
			updateDerivedColumnCmd(),
			deleteDerivedColumnCmd(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func newClient(cmd *cli.Command) *api.Client {
	return api.NewClient(cmd.String("api-key"))
}

func authV2Cmd() *cli.Command {
	return &cli.Command{
		Name:  "auth-v2",
		Usage: "Show management API key info and permissions (v2)",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			auth, err := client.GetAuthV2(ctx)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(auth)
		},
	}
}

func listBoardsCmd() *cli.Command {
	return &cli.Command{
		Name:  "boards",
		Usage: "List all boards",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			boards, err := client.ListBoards(ctx)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(boards)
		},
	}
}

func getBoardCmd() *cli.Command {
	return &cli.Command{
		Name:  "get-board",
		Usage: "Get a board by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Board ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			board, err := client.GetBoard(ctx, cmd.String("id"))
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(board)
		},
	}
}

func createBoardCmd() *cli.Command {
	return &cli.Command{
		Name:  "create-board",
		Usage: "Create a new board",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Board name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Board description",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			board := &api.Board{
				Name:        cmd.String("name"),
				Description: cmd.String("description"),
				Type:        "flexible",
			}

			created, err := client.CreateBoard(ctx, board)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(created)
		},
	}
}

func updateBoardCmd() *cli.Command {
	return &cli.Command{
		Name:  "update-board",
		Usage: "Update a board by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Board ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Board name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Board description",
			},
			&cli.StringFlag{
				Name:  "query-id",
				Usage: "Query ID for a panel",
			},
			&cli.StringFlag{
				Name:  "query-annotation-id",
				Usage: "Query annotation ID for a panel",
			},
			&cli.StringFlag{
				Name:  "query-style",
				Usage: "Query style (graph, table, combo)",
				Value: "graph",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			board := &api.Board{
				Name:        cmd.String("name"),
				Description: cmd.String("description"),
				Type:        "flexible",
			}

			if qid := cmd.String("query-id"); qid != "" {
				style := cmd.String("query-style")
				if style == "" {
					style = "graph"
				}
				board.Panels = []api.BoardPanel{
					{
						Type: "query",
						QueryPanel: &api.QueryPanel{
							QueryID:           qid,
							QueryAnnotationID: cmd.String("query-annotation-id"),
							QueryStyle:        style,
						},
					},
				}
			}

			updated, err := client.UpdateBoard(ctx, cmd.String("id"), board)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(updated)
		},
	}
}

func deleteBoardCmd() *cli.Command {
	return &cli.Command{
		Name:  "delete-board",
		Usage: "Delete a board by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Board ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			return client.DeleteBoard(ctx, cmd.String("id"))
		},
	}
}

func listBoardViewsCmd() *cli.Command {
	return &cli.Command{
		Name:  "board-views",
		Usage: "List all views for a board",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "board-id",
				Usage:    "Board ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			views, err := client.ListBoardViews(ctx, cmd.String("board-id"))
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(views)
		},
	}
}

func getBoardViewCmd() *cli.Command {
	return &cli.Command{
		Name:  "get-board-view",
		Usage: "Get a board view by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "board-id",
				Usage:    "Board ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "view-id",
				Usage:    "View ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			view, err := client.GetBoardView(ctx, cmd.String("board-id"), cmd.String("view-id"))
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(view)
		},
	}
}

func createBoardViewCmd() *cli.Command {
	return &cli.Command{
		Name:  "create-board-view",
		Usage: "Create a new view for a board",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "board-id",
				Usage:    "Board ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "View name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "filter-column",
				Usage:    "Filter column name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "filter-op",
				Usage:    "Filter operation (e.g. =, !=, >, <, starts-with)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "filter-value",
				Usage: "Filter value",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			filter := api.BoardViewFilter{
				Column:    cmd.String("filter-column"),
				Operation: cmd.String("filter-op"),
			}
			if v := cmd.String("filter-value"); v != "" {
				filter.Value = v
			}

			view := &api.BoardView{
				Name:    cmd.String("name"),
				Filters: []api.BoardViewFilter{filter},
			}

			created, err := client.CreateBoardView(ctx, cmd.String("board-id"), view)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(created)
		},
	}
}

func updateBoardViewCmd() *cli.Command {
	return &cli.Command{
		Name:  "update-board-view",
		Usage: "Update a board view by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "board-id",
				Usage:    "Board ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "view-id",
				Usage:    "View ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "View name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "filter-column",
				Usage:    "Filter column name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "filter-op",
				Usage:    "Filter operation (e.g. =, !=, >, <, starts-with)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "filter-value",
				Usage: "Filter value",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			filter := api.BoardViewFilter{
				Column:    cmd.String("filter-column"),
				Operation: cmd.String("filter-op"),
			}
			if v := cmd.String("filter-value"); v != "" {
				filter.Value = v
			}

			view := &api.BoardView{
				Name:    cmd.String("name"),
				Filters: []api.BoardViewFilter{filter},
			}

			updated, err := client.UpdateBoardView(ctx, cmd.String("board-id"), cmd.String("view-id"), view)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(updated)
		},
	}
}

func deleteBoardViewCmd() *cli.Command {
	return &cli.Command{
		Name:  "delete-board-view",
		Usage: "Delete a board view by ID",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "board-id",
				Usage:    "Board ID",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "view-id",
				Usage:    "View ID",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)
			return client.DeleteBoardView(ctx, cmd.String("board-id"), cmd.String("view-id"))
		},
	}
}

func getQueryCmd() *cli.Command {
	return &cli.Command{
		Name:  "get-query",
		Usage: "Get a query by ID",
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

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(query)
		},
	}
}

func createQueryCmd() *cli.Command {
	return &cli.Command{
		Name:  "create-query",
		Usage: "Create a new query",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "calculation-op",
				Usage:    "Calculation operation (e.g. COUNT, AVG, P99)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "calculation-column",
				Usage: "Calculation column (not needed for COUNT)",
			},
			&cli.StringFlag{
				Name:  "breakdown",
				Usage: "Breakdown column",
			},
			&cli.StringFlag{
				Name:  "filter-column",
				Usage: "Filter column",
			},
			&cli.StringFlag{
				Name:  "filter-op",
				Usage: "Filter operation",
			},
			&cli.StringFlag{
				Name:  "filter-value",
				Usage: "Filter value",
			},
			&cli.IntFlag{
				Name:  "time-range",
				Usage: "Time range in seconds",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			calc := api.Calculation{
				Op: cmd.String("calculation-op"),
			}
			if v := cmd.String("calculation-column"); v != "" {
				calc.Column = v
			}

			query := &api.Query{
				Calculations: []api.Calculation{calc},
			}

			if v := cmd.String("breakdown"); v != "" {
				query.Breakdowns = []string{v}
			}

			if col := cmd.String("filter-column"); col != "" {
				f := api.QueryFilter{
					Column: col,
					Op:     cmd.String("filter-op"),
				}
				if v := cmd.String("filter-value"); v != "" {
					f.Value = v
				}
				query.Filters = []api.QueryFilter{f}
			}

			if v := cmd.Int("time-range"); v != 0 {
				query.TimeRange = int(v)
			}

			created, err := client.CreateQuery(ctx, cmd.String("dataset"), query)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(created)
		},
	}
}

func createQueryAnnotationCmd() *cli.Command {
	return &cli.Command{
		Name:  "create-query-annotation",
		Usage: "Create a query annotation",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dataset",
				Usage:    "Dataset slug (use __all__ for environment-wide)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "query-id",
				Usage:    "Query ID to annotate",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Usage:    "Annotation name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "description",
				Usage: "Annotation description",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			annotation := &api.QueryAnnotation{
				Name:    cmd.String("name"),
				QueryID: cmd.String("query-id"),
			}
			if v := cmd.String("description"); v != "" {
				annotation.Description = v
			}

			created, err := client.CreateQueryAnnotation(ctx, cmd.String("dataset"), annotation)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(created)
		},
	}
}

func listDerivedColumnsCmd() *cli.Command {
	return &cli.Command{
		Name:  "derived-columns",
		Usage: "List all calculated fields (derived columns)",
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

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(cols)
		},
	}
}

func getDerivedColumnCmd() *cli.Command {
	return &cli.Command{
		Name:  "get-derived-column",
		Usage: "Get a calculated field (derived column) by ID",
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

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(col)
		},
	}
}

func createDerivedColumnCmd() *cli.Command {
	return &cli.Command{
		Name:  "create-derived-column",
		Usage: "Create a new calculated field (derived column)",
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

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(created)
		},
	}
}

func updateDerivedColumnCmd() *cli.Command {
	return &cli.Command{
		Name:  "update-derived-column",
		Usage: "Update a calculated field (derived column) by ID",
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

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(updated)
		},
	}
}

func deleteDerivedColumnCmd() *cli.Command {
	return &cli.Command{
		Name:  "delete-derived-column",
		Usage: "Delete a calculated field (derived column) by ID",
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

func authCmd() *cli.Command {
	return &cli.Command{
		Name:  "auth",
		Usage: "Show API key info and permissions",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			auth, err := client.GetAuth(ctx)
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(auth)
		},
	}
}
