package main_test

import (
	"testing"
)

func TestCreateAndGetQueryResultCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	// First, create a query definition.
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
	queryID, ok := query["id"].(string)
	if !ok || queryID == "" {
		t.Fatal("expected non-empty query id")
	}

	// Execute the query and wait for results.
	stdout, stderr, exitCode := runCLIWithKey(t,
		"create-query-result",
		"--dataset", dataset,
		"--query-id", queryID,
		"--timeout", "30",
	)
	if exitCode != 0 {
		t.Fatalf("create-query-result failed with exit code %d\nstdout: %s\nstderr: %s", exitCode, stdout, stderr)
	}

	result := parseJSON(t, stdout)

	resultID, ok := result["id"].(string)
	if !ok || resultID == "" {
		t.Fatal("expected non-empty query result id")
	}

	complete, ok := result["complete"].(bool)
	if !ok || !complete {
		t.Fatalf("expected complete=true, got %v", result["complete"])
	}

	data, ok := result["data"].(map[string]any)
	if !ok {
		t.Fatal("expected data field in query result")
	}

	if _, ok := data["results"]; !ok {
		t.Error("expected results field in query result data")
	}

	// Fetch the same result again by ID.
	stdout, _, exitCode = runCLIWithKey(t,
		"get-query-result",
		"--dataset", dataset,
		"--id", resultID,
	)
	if exitCode != 0 {
		t.Fatalf("get-query-result failed with exit code %d", exitCode)
	}

	got := parseJSON(t, stdout)
	if got["id"] != resultID {
		t.Errorf("expected result id %q, got %v", resultID, got["id"])
	}
}
