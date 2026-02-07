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
			createBoardCmd(),
			deleteBoardCmd(),
			createBoardViewCmd(),
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
