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

func TestCreateQueryMultipleFiltersCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	stdout, stderr, exitCode := runCLIWithKey(t,
		"create-query",
		"--dataset", dataset,
		"--calculation-op", "COUNT",
		"--filter", "service.name exists",
		"--filter", "service.instance.id exists",
		"--time-range", "7200",
	)
	if exitCode != 0 {
		t.Fatalf("create-query with multiple filters failed with exit code %d: %s", exitCode, stderr)
	}

	query := parseJSON(t, stdout)
	filters, ok := query["filters"].([]any)
	if !ok || len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %v", query["filters"])
	}

	f0 := filters[0].(map[string]any)
	if f0["column"] != "service.name" || f0["op"] != "exists" {
		t.Errorf("unexpected first filter: %v", f0)
	}
	f1 := filters[1].(map[string]any)
	if f1["column"] != "service.instance.id" || f1["op"] != "exists" {
		t.Errorf("unexpected second filter: %v", f1)
	}
}

func TestCreateQueryMultipleBreakdownsCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	stdout, stderr, exitCode := runCLIWithKey(t,
		"create-query",
		"--dataset", dataset,
		"--calculation-op", "COUNT",
		"--breakdown", "service.name",
		"--breakdown", "service.instance.id",
		"--time-range", "7200",
	)
	if exitCode != 0 {
		t.Fatalf("create-query with multiple breakdowns failed with exit code %d: %s", exitCode, stderr)
	}

	query := parseJSON(t, stdout)
	breakdowns, ok := query["breakdowns"].([]any)
	if !ok || len(breakdowns) != 2 {
		t.Fatalf("expected 2 breakdowns, got %v", query["breakdowns"])
	}
	if breakdowns[0] != "service.name" {
		t.Errorf("expected first breakdown 'service.name', got %v", breakdowns[0])
	}
	if breakdowns[1] != "service.instance.id" {
		t.Errorf("expected second breakdown 'service.instance.id', got %v", breakdowns[1])
	}
}

func TestCreateQueryInvalidFilterCLI(t *testing.T) {
	_, stderr, exitCode := runCLI(t,
		"--api-key", "fake-key",
		"create-query",
		"--dataset", "test",
		"--calculation-op", "COUNT",
		"--filter", "just-a-column",
	)
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code for invalid filter")
	}
	if !strings.Contains(stderr, "invalid filter") {
		t.Errorf("expected error about invalid filter, got: %s", stderr)
	}
}

func TestCreateQueryFilterMissingValueCLI(t *testing.T) {
	_, stderr, exitCode := runCLI(t,
		"--api-key", "fake-key",
		"create-query",
		"--dataset", "test",
		"--calculation-op", "COUNT",
		"--filter", "col1 =",
	)
	if exitCode == 0 {
		t.Fatal("expected non-zero exit code for filter missing value")
	}
	if !strings.Contains(stderr, "requires a value") {
		t.Errorf("expected error about missing value, got: %s", stderr)
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
