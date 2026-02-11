package main_test

import (
	"strings"
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

func TestCreateQueryMultipleCalculationsCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-query",
		"--dataset", dataset,
		"--calculation-op", "COUNT",
		"--calculation-op", "CONCURRENCY",
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
	if !ok || len(calcs) != 2 {
		t.Fatalf("expected 2 calculations, got %v", query["calculations"])
	}

	calc0 := calcs[0].(map[string]any)
	if calc0["op"] != "COUNT" {
		t.Errorf("expected first calculation op 'COUNT', got %v", calc0["op"])
	}

	calc1 := calcs[1].(map[string]any)
	if calc1["op"] != "CONCURRENCY" {
		t.Errorf("expected second calculation op 'CONCURRENCY', got %v", calc1["op"])
	}
}

func TestCreateQueryMismatchedCalculationsCLI(t *testing.T) {
	_, stderr, exitCode := runCLI(t,
		"--api-key", "fake-key",
		"create-query",
		"--dataset", "test",
		"--calculation-op", "COUNT",
		"--calculation-op", "AVG",
		"--calculation-column", "only_one",
	)
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code for mismatched calculation flags")
	}
	if !strings.Contains(stderr, "must match") {
		t.Errorf("expected error about mismatched counts, got: %s", stderr)
	}
}
