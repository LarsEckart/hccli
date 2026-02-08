package api_test

import (
	"context"
	"os"
	"testing"

	"github.com/LarsEckart/hccli/api"
)

func TestDerivedColumnCRUD_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}
	if os.Getenv("HONEYCOMB_DATASET") == "" {
		t.Skip("HONEYCOMB_DATASET not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)
	dataset := testDataset(t)

	col, err := client.CreateDerivedColumn(context.Background(), dataset, &api.DerivedColumn{
		Alias:       "test_derived_col",
		Expression:  "INT(1)",
		Description: "test derived column",
	})
	if err != nil {
		t.Fatalf("CreateDerivedColumn failed: %v", err)
	}
	if col.ID == "" {
		t.Fatal("expected non-empty derived column ID")
	}
	if col.Alias != "test_derived_col" {
		t.Errorf("expected alias %q, got %q", "test_derived_col", col.Alias)
	}

	got, err := client.GetDerivedColumn(context.Background(), dataset, col.ID)
	if err != nil {
		t.Fatalf("GetDerivedColumn failed: %v", err)
	}
	if got.ID != col.ID {
		t.Errorf("expected ID %q, got %q", col.ID, got.ID)
	}

	updated, err := client.UpdateDerivedColumn(context.Background(), dataset, col.ID, &api.DerivedColumn{
		Alias:       "test_derived_col_updated",
		Expression:  "INT(2)",
		Description: "updated description",
	})
	if err != nil {
		t.Fatalf("UpdateDerivedColumn failed: %v", err)
	}
	if updated.Alias != "test_derived_col_updated" {
		t.Errorf("expected alias %q, got %q", "test_derived_col_updated", updated.Alias)
	}

	cols, err := client.ListDerivedColumns(context.Background(), dataset)
	if err != nil {
		t.Fatalf("ListDerivedColumns failed: %v", err)
	}
	found := false
	for _, c := range cols {
		if c.ID == col.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find derived column %q in list", col.ID)
	}

	if err := client.DeleteDerivedColumn(context.Background(), dataset, col.ID); err != nil {
		t.Fatalf("DeleteDerivedColumn failed: %v", err)
	}
}
