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

func TestGetAuthV2_Smoke(t *testing.T) {
	keyID := os.Getenv("HONEYCOMB_MANAGEMENT_API_KEY_ID")
	keySecret := os.Getenv("HONEYCOMB_MANAGEMENT_API_KEY_SECRET")
	if keyID == "" || keySecret == "" {
		t.Skip("HONEYCOMB_MANAGEMENT_API_KEY_ID or HONEYCOMB_MANAGEMENT_API_KEY_SECRET not set, skipping smoke test")
	}

	client := api.NewClient(keyID + ":" + keySecret)
	auth, err := client.GetAuthV2(context.Background())
	if err != nil {
		t.Fatalf("GetAuthV2 failed: %v", err)
	}

	if auth.Data.ID == "" {
		t.Error("expected non-empty Data.ID")
	}
	if auth.Data.Type != "api-keys" {
		t.Errorf("expected Data.Type 'api-keys', got %q", auth.Data.Type)
	}
}
