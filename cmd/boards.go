package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/LarsEckart/hccli/api"
	"github.com/urfave/cli/v3"
)

func ListBoardsCmd() *cli.Command {
	return &cli.Command{
		Name:     "boards",
		Category: "Boards",
		Usage:    "List all boards",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			boards, err := client.ListBoards(ctx)
			if err != nil {
				return err
			}

			return printJSON(boards)
		},
	}
}

func GetBoardCmd() *cli.Command {
	return &cli.Command{
		Name:     "get-board",
		Category: "Boards",
		Usage:    "Get a board by ID",
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

			return printJSON(board)
		},
	}
}

func CreateBoardCmd() *cli.Command {
	return &cli.Command{
		Name:     "create-board",
		Category: "Boards",
		Usage:    "Create a new board",
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

			return printJSON(created)
		},
	}
}

func UpdateBoardCmd() *cli.Command {
	return &cli.Command{
		Name:     "update-board",
		Category: "Boards",
		Usage:    "Update a board by ID",
		Description: `Update a board's name, description, and panels.

To replace the full set of panels, use --panels-json with a JSON array
from get-board output. This enables adding or removing individual panels:

  # Get current panels, remove index 1, and update
  PANELS=$(hccli get-board --id ID | jq 'del(.panels[1]) | .panels')
  hccli update-board --id ID --name "my board" --panels-json "$PANELS"`,
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
			&cli.StringFlag{
				Name:  "panels-json",
				Usage: "Full JSON array of board panels; use get-board output to build it (overrides --query-id)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(cmd)

			board := &api.Board{
				Name:        cmd.String("name"),
				Description: cmd.String("description"),
				Type:        "flexible",
			}

			if pj := cmd.String("panels-json"); pj != "" {
				var panels []api.BoardPanel
				if err := json.Unmarshal([]byte(pj), &panels); err != nil {
					return fmt.Errorf("parsing panels-json: %w", err)
				}
				board.Panels = panels
			} else if qid := cmd.String("query-id"); qid != "" {
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

			return printJSON(updated)
		},
	}
}

func DeleteBoardCmd() *cli.Command {
	return &cli.Command{
		Name:     "delete-board",
		Category: "Boards",
		Usage:    "Delete a board by ID",
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
