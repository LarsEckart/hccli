package api_test

import (
	"context"
	"os"
	"testing"

	"github.com/LarsEckart/hccli/api"
)

func testDataset(t *testing.T) string {
	t.Helper()
	ds := os.Getenv("HONEYCOMB_DATASET")
	if ds == "" {
		ds = "__all__"
	}
	return ds
}

func TestCreateAndGetQuery_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}
	if os.Getenv("HONEYCOMB_DATASET") == "" {
		t.Skip("HONEYCOMB_DATASET not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)
	dataset := testDataset(t)

	query, err := client.CreateQuery(context.Background(), dataset, &api.Query{
		Calculations: []api.Calculation{
			{Op: "COUNT"},
		},
		TimeRange: 7200,
	})
	if err != nil {
		t.Fatalf("CreateQuery failed: %v", err)
	}

	if query.ID == "" {
		t.Fatal("expected non-empty query ID")
	}
	if len(query.Calculations) != 1 {
		t.Fatalf("expected 1 calculation, got %d", len(query.Calculations))
	}
	if query.Calculations[0].Op != "COUNT" {
		t.Errorf("expected calculation op 'COUNT', got %q", query.Calculations[0].Op)
	}

	got, err := client.GetQuery(context.Background(), dataset, query.ID)
	if err != nil {
		t.Fatalf("GetQuery failed: %v", err)
	}

	if got.ID != query.ID {
		t.Errorf("expected query ID %q, got %q", query.ID, got.ID)
	}
}
