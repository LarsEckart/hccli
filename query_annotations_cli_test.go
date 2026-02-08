package main_test

import (
	"testing"
)

func TestQueryAnnotationCRUDCLI_Smoke(t *testing.T) {
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
	id, ok := annotation["id"].(string)
	if !ok || id == "" {
		t.Fatal("expected non-empty annotation id")
	}
	if annotation["name"] != "hccli smoke test annotation" {
		t.Errorf("expected name 'hccli smoke test annotation', got %v", annotation["name"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "get-query-annotation", "--dataset", dataset, "--id", id)
	if exitCode != 0 {
		t.Fatalf("get-query-annotation failed with exit code %d", exitCode)
	}
	got := parseJSON(t, stdout)
	if got["id"] != id {
		t.Errorf("expected id %q, got %v", id, got["id"])
	}

	stdout, _, exitCode = runCLIWithKey(t,
		"update-query-annotation",
		"--dataset", dataset,
		"--id", id,
		"--name", "updated annotation",
		"--query-id", queryID,
		"--description", "updated description",
	)
	if exitCode != 0 {
		t.Fatalf("update-query-annotation failed with exit code %d", exitCode)
	}
	updated := parseJSON(t, stdout)
	if updated["name"] != "updated annotation" {
		t.Errorf("expected name 'updated annotation', got %v", updated["name"])
	}
	if updated["description"] != "updated description" {
		t.Errorf("expected description 'updated description', got %v", updated["description"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "query-annotations", "--dataset", dataset)
	if exitCode != 0 {
		t.Fatalf("query-annotations failed with exit code %d", exitCode)
	}
	arr := parseJSONArray(t, stdout)
	found := false
	for _, item := range arr {
		a := item.(map[string]any)
		if a["id"] == id {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find annotation %q in list", id)
	}

	_, _, exitCode = runCLIWithKey(t, "delete-query-annotation", "--dataset", dataset, "--id", id)
	if exitCode != 0 {
		t.Fatalf("delete-query-annotation failed with exit code %d", exitCode)
	}
}
