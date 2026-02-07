package api_test

import (
	"context"
	"os"
	"testing"

	"github.com/LarsEckart/hccli/api"
)

func TestListBoards_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)
	boards, err := client.ListBoards(context.Background())
	if err != nil {
		t.Fatalf("ListBoards failed: %v", err)
	}

	if boards == nil {
		t.Error("expected non-nil boards slice")
	}
}

func TestCreateAndDeleteBoard_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)
	board, err := client.CreateBoard(context.Background(), &api.Board{
		Name: "hccli smoke test board",
		Type: "flexible",
	})
	if err != nil {
		t.Fatalf("CreateBoard failed: %v", err)
	}

	if board.ID == "" {
		t.Fatal("expected non-empty board ID")
	}
	if board.Name != "hccli smoke test board" {
		t.Errorf("expected name 'hccli smoke test board', got %q", board.Name)
	}

	err = client.DeleteBoard(context.Background(), board.ID)
	if err != nil {
		t.Fatalf("DeleteBoard failed: %v", err)
	}
}

func TestCreateBoardView_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)

	board, err := client.CreateBoard(context.Background(), &api.Board{
		Name: "hccli board view smoke test",
		Type: "flexible",
	})
	if err != nil {
		t.Fatalf("CreateBoard failed: %v", err)
	}
	defer func() {
		_ = client.DeleteBoard(context.Background(), board.ID)
	}()

	view, err := client.CreateBoardView(context.Background(), board.ID, &api.BoardView{
		Name: "test view",
		Filters: []api.BoardViewFilter{
			{
				Column:    "status",
				Operation: "=",
				Value:     "200",
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateBoardView failed: %v", err)
	}

	if view.ID == "" {
		t.Fatal("expected non-empty view ID")
	}
	if view.Name != "test view" {
		t.Errorf("expected name 'test view', got %q", view.Name)
	}
	if len(view.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(view.Filters))
	}
	if view.Filters[0].Column != "status" {
		t.Errorf("expected filter column 'status', got %q", view.Filters[0].Column)
	}
}
