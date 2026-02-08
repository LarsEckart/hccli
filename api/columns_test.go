package api_test

import (
	"context"
	"os"
	"testing"

	"github.com/LarsEckart/hccli/api"
)

func TestColumnCRUD_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}
	if os.Getenv("HONEYCOMB_DATASET") == "" {
		t.Skip("HONEYCOMB_DATASET not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)
	dataset := testDataset(t)

	hidden := false
	col, err := client.CreateColumn(context.Background(), dataset, &api.Column{
		KeyName:     "test_column",
		Type:        "integer",
		Description: "test column",
		Hidden:      &hidden,
	})
	if err != nil {
		t.Fatalf("CreateColumn failed: %v", err)
	}
	if col.ID == "" {
		t.Fatal("expected non-empty column ID")
	}
	if col.KeyName != "test_column" {
		t.Errorf("expected key_name %q, got %q", "test_column", col.KeyName)
	}

	got, err := client.GetColumn(context.Background(), dataset, col.ID)
	if err != nil {
		t.Fatalf("GetColumn failed: %v", err)
	}
	if got.ID != col.ID {
		t.Errorf("expected ID %q, got %q", col.ID, got.ID)
	}

	updated, err := client.UpdateColumn(context.Background(), dataset, col.ID, &api.Column{
		Description: "updated description",
		Type:        "integer",
	})
	if err != nil {
		t.Fatalf("UpdateColumn failed: %v", err)
	}
	if updated.Description != "updated description" {
		t.Errorf("expected description %q, got %q", "updated description", updated.Description)
	}

	cols, err := client.ListColumns(context.Background(), dataset)
	if err != nil {
		t.Fatalf("ListColumns failed: %v", err)
	}
	found := false
	for _, c := range cols {
		if c.ID == col.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find column %q in list", col.ID)
	}

	if err := client.DeleteColumn(context.Background(), dataset, col.ID); err != nil {
		t.Fatalf("DeleteColumn failed: %v", err)
	}
}
