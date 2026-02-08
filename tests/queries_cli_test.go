package main_test

import (
	"testing"
)

func TestCreateAndGetQueryCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-query",
		"--dataset", dataset,
		"--calculation-op", "COUNT",
		"--time-range", "7200",
	)
	if exitCode != 0 {
		t.Fatalf("create-query failed with exit code %d", exitCode)
	}

	query := parseJSON(t, stdout)
	id, ok := query["id"].(string)
	if !ok || id == "" {
		t.Fatal("expected non-empty query id")
	}

	calcs, ok := query["calculations"].([]any)
	if !ok || len(calcs) != 1 {
		t.Fatalf("expected 1 calculation, got %v", query["calculations"])
	}
	calc := calcs[0].(map[string]any)
	if calc["op"] != "COUNT" {
		t.Errorf("expected calculation op 'COUNT', got %v", calc["op"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "get-query", "--dataset", dataset, "--id", id)
	if exitCode != 0 {
		t.Fatalf("get-query failed with exit code %d", exitCode)
	}

	got := parseJSON(t, stdout)
	if got["id"] != id {
		t.Errorf("expected query id %q, got %v", id, got["id"])
	}
}
