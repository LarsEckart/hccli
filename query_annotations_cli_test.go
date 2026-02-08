package main_test

import (
	"testing"
)

func TestCreateQueryAnnotationCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-query",
		"--dataset", dataset,
		"--calculation-op", "COUNT",
	)
	if exitCode != 0 {
		t.Fatalf("create-query failed with exit code %d", exitCode)
	}
	query := parseJSON(t, stdout)
	queryID := query["id"].(string)

	stdout, _, exitCode = runCLIWithKey(t,
		"create-query-annotation",
		"--dataset", dataset,
		"--query-id", queryID,
		"--name", "hccli smoke test annotation",
	)
	if exitCode != 0 {
		t.Fatalf("create-query-annotation failed with exit code %d", exitCode)
	}

	annotation := parseJSON(t, stdout)
	if annotation["id"] == nil || annotation["id"] == "" {
		t.Fatal("expected non-empty annotation id")
	}
	if annotation["name"] != "hccli smoke test annotation" {
		t.Errorf("expected name 'hccli smoke test annotation', got %v", annotation["name"])
	}
}
