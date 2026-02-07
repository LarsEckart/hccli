package api_test

import (
	"context"
	"os"
	"testing"

	"github.com/LarsEckart/hccli/api"
)

func TestGetAuth_Smoke(t *testing.T) {
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}

	client := api.NewClient(apiKey)
	auth, err := client.GetAuth(context.Background())
	if err != nil {
		t.Fatalf("GetAuth failed: %v", err)
	}

	if auth.ID == "" {
		t.Error("expected non-empty ID")
	}
	if auth.Type == "" {
		t.Error("expected non-empty Type")
	}
	if auth.Team.Slug == "" {
		t.Error("expected non-empty Team slug")
	}
}
