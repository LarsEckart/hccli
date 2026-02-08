package cmd

import (
	"context"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func ListBoardViewsCmd() *cli.Command {
	return &cli.Command{
		Name:     "board-views",
		Category: "Board Views",
		Usage:    "List all views for a board",
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

			return printJSON(views)
		},
	}
}

func GetBoardViewCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-board-view",
		Category: "Board Views",
		Usage:    "Get a board view by ID",
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

			return printJSON(view)
		},
	}
}

func CreateBoardViewCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-board-view",
		Category: "Board Views",
		Usage:    "Create a new view for a board",
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

			return printJSON(created)
		},
	}
}

func UpdateBoardViewCmd() *cli.Command {
	return &cli.Command{
		Name:     "update-board-view",
		Category: "Board Views",
		Usage:    "Update a board view by ID",
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

			return printJSON(updated)
		},
	}
}

func DeleteBoardViewCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete-board-view",
		Category: "Board Views",
		Usage:    "Delete a board view by ID",
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
