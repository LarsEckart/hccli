package api_test

import (
	"context"
	"os"
	"testing"

	"github.com/LarsEckart/hccli/api"
)

func TestCreateQueryAnnotation_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)
	ctx := context.Background()

	query, err := client.CreateQuery(ctx, "__all__", &api.Query{
		Calculations: []api.Calculation{{Op: "COUNT"}},
	})
	if err != nil {
		t.Fatalf("CreateQuery failed: %v", err)
	}

	annotation, err := client.CreateQueryAnnotation(ctx, "__all__", &api.QueryAnnotation{
		Name:    "hccli smoke test annotation",
		QueryID: query.ID,
	})
	if err != nil {
		t.Fatalf("CreateQueryAnnotation failed: %v", err)
	}

	if annotation.ID == "" {
		t.Fatal("expected non-empty annotation ID")
	}
	if annotation.Name != "hccli smoke test annotation" {
		t.Errorf("expected name 'hccli smoke test annotation', got %q", annotation.Name)
	}
}
