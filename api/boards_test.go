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

func TestUpdateBoard_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)

	board, err := client.CreateBoard(context.Background(), &api.Board{
		Name: "hccli update smoke test",
		Type: "flexible",
	})
	if err != nil {
		t.Fatalf("CreateBoard failed: %v", err)
	}
	defer func() {
		_ = client.DeleteBoard(context.Background(), board.ID)
	}()

	board.Name = "hccli update smoke test updated"
	board.Description = "updated description"

	updated, err := client.UpdateBoard(context.Background(), board.ID, board)
	if err != nil {
		t.Fatalf("UpdateBoard failed: %v", err)
	}

	if updated.Name != "hccli update smoke test updated" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}
	if updated.Description != "updated description" {
		t.Errorf("expected updated description, got %q", updated.Description)
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

func TestListBoardViews_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)

	board, err := client.CreateBoard(context.Background(), &api.Board{
		Name: "hccli list views smoke test",
		Type: "flexible",
	})
	if err != nil {
		t.Fatalf("CreateBoard failed: %v", err)
	}
	defer func() {
		_ = client.DeleteBoard(context.Background(), board.ID)
	}()

	_, err = client.CreateBoardView(context.Background(), board.ID, &api.BoardView{
		Name: "view1",
		Filters: []api.BoardViewFilter{
			{Column: "status", Operation: "=", Value: "200"},
		},
	})
	if err != nil {
		t.Fatalf("CreateBoardView failed: %v", err)
	}

	views, err := client.ListBoardViews(context.Background(), board.ID)
	if err != nil {
		t.Fatalf("ListBoardViews failed: %v", err)
	}

	if len(views) < 1 {
		t.Fatalf("expected at least 1 view, got %d", len(views))
	}
}

func TestGetBoardView_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)

	board, err := client.CreateBoard(context.Background(), &api.Board{
		Name: "hccli get view smoke test",
		Type: "flexible",
	})
	if err != nil {
		t.Fatalf("CreateBoard failed: %v", err)
	}
	defer func() {
		_ = client.DeleteBoard(context.Background(), board.ID)
	}()

	created, err := client.CreateBoardView(context.Background(), board.ID, &api.BoardView{
		Name: "view to get",
		Filters: []api.BoardViewFilter{
			{Column: "status", Operation: "=", Value: "200"},
		},
	})
	if err != nil {
		t.Fatalf("CreateBoardView failed: %v", err)
	}

	view, err := client.GetBoardView(context.Background(), board.ID, created.ID)
	if err != nil {
		t.Fatalf("GetBoardView failed: %v", err)
	}

	if view.ID != created.ID {
		t.Errorf("expected view ID %q, got %q", created.ID, view.ID)
	}
	if view.Name != "view to get" {
		t.Errorf("expected name 'view to get', got %q", view.Name)
	}
}

func TestUpdateBoardView_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)

	board, err := client.CreateBoard(context.Background(), &api.Board{
		Name: "hccli update view smoke test",
		Type: "flexible",
	})
	if err != nil {
		t.Fatalf("CreateBoard failed: %v", err)
	}
	defer func() {
		_ = client.DeleteBoard(context.Background(), board.ID)
	}()

	created, err := client.CreateBoardView(context.Background(), board.ID, &api.BoardView{
		Name: "original view",
		Filters: []api.BoardViewFilter{
			{Column: "status", Operation: "=", Value: "200"},
		},
	})
	if err != nil {
		t.Fatalf("CreateBoardView failed: %v", err)
	}

	updated, err := client.UpdateBoardView(context.Background(), board.ID, created.ID, &api.BoardView{
		Name: "updated view",
		Filters: []api.BoardViewFilter{
			{Column: "status", Operation: "!=", Value: "500"},
		},
	})
	if err != nil {
		t.Fatalf("UpdateBoardView failed: %v", err)
	}

	if updated.Name != "updated view" {
		t.Errorf("expected name 'updated view', got %q", updated.Name)
	}
	if len(updated.Filters) != 1 {
		t.Fatalf("expected 1 filter, got %d", len(updated.Filters))
	}
	if updated.Filters[0].Operation != "!=" {
		t.Errorf("expected filter op '!=', got %q", updated.Filters[0].Operation)
	}
}

func TestDeleteBoardView_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)

	board, err := client.CreateBoard(context.Background(), &api.Board{
		Name: "hccli delete view smoke test",
		Type: "flexible",
	})
	if err != nil {
		t.Fatalf("CreateBoard failed: %v", err)
	}
	defer func() {
		_ = client.DeleteBoard(context.Background(), board.ID)
	}()

	created, err := client.CreateBoardView(context.Background(), board.ID, &api.BoardView{
		Name: "view to delete",
		Filters: []api.BoardViewFilter{
			{Column: "status", Operation: "=", Value: "200"},
		},
	})
	if err != nil {
		t.Fatalf("CreateBoardView failed: %v", err)
	}

	err = client.DeleteBoardView(context.Background(), board.ID, created.ID)
	if err != nil {
		t.Fatalf("DeleteBoardView failed: %v", err)
	}

	_, err = client.GetBoardView(context.Background(), board.ID, created.ID)
	if err == nil {
		t.Error("expected error when getting deleted view, got nil")
	}
}
