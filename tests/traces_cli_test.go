package main_test

import (
	"testing"
)

func TestGetTraceCLI_Smoke(t *testing.T) {
	stdout, _, exitCode := runCLIWithKey(t, "get-trace", "--trace-id", "abc123", "--dataset", "frontend")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	m := parseJSON(t, stdout)
	if m["trace_id"] != "abc123" {
		t.Errorf("expected trace_id 'abc123', got %v", m["trace_id"])
	}
	if m["dataset"] != "frontend" {
		t.Errorf("expected dataset 'frontend', got %v", m["dataset"])
	}
	if m["url"] == nil || m["url"] == "" {
		t.Error("expected non-empty url")
	}
	if m["team"] == nil || m["team"] == "" {
		t.Error("expected non-empty team")
	}
	if m["environment"] == nil || m["environment"] == "" {
		t.Error("expected non-empty environment")
	}
}

func TestGetTraceCLI_MissingTraceID(t *testing.T) {
	_, _, exitCode := runCLIWithKey(t, "get-trace", "--dataset", "frontend")
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code when trace-id is missing")
	}
}

func TestGetTraceCLI_MissingDataset(t *testing.T) {
	_, _, exitCode := runCLIWithKey(t, "get-trace", "--trace-id", "abc123")
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code when dataset is missing")
	}
}
