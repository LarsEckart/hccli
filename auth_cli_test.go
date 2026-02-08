package main_test

import (
	"os"
	"testing"
)

func TestAuthCLI_Smoke(t *testing.T) {
	stdout, _, exitCode := runCLIWithKey(t, "auth")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	m := parseJSON(t, stdout)
	if m["id"] == nil || m["id"] == "" {
		t.Error("expected non-empty id")
	}
	if m["type"] == nil || m["type"] == "" {
		t.Error("expected non-empty type")
	}
	team, ok := m["team"].(map[string]any)
	if !ok {
		t.Fatal("expected team object")
	}
	if team["slug"] == nil || team["slug"] == "" {
		t.Error("expected non-empty team slug")
	}
}

func TestAuthV2CLI_Smoke(t *testing.T) {
	keyID := os.Getenv("HONEYCOMB_MANAGEMENT_API_KEY_ID")
	keySecret := os.Getenv("HONEYCOMB_MANAGEMENT_API_KEY_SECRET")
	if keyID == "" || keySecret == "" {
		t.Skip("HONEYCOMB_MANAGEMENT_API_KEY_ID or HONEYCOMB_MANAGEMENT_API_KEY_SECRET not set, skipping smoke test")
	}

	stdout, _, exitCode := runCLI(t, "--api-key", keyID+":"+keySecret, "auth-v2")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	m := parseJSON(t, stdout)
	data, ok := m["data"].(map[string]any)
	if !ok {
		t.Fatal("expected data object")
	}
	if data["id"] == nil || data["id"] == "" {
		t.Error("expected non-empty data.id")
	}
	if data["type"] != "api-keys" {
		t.Errorf("expected data.type 'api-keys', got %v", data["type"])
	}
}
