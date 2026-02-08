package main_test

import (
	"testing"
)

func TestColumnCRUDCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-column",
		"--dataset", dataset,
		"--key-name", "test_column",
		"--type", "integer",
		"--description", "test column",
	)
	if exitCode != 0 {
		t.Fatalf("create-column failed with exit code %d", exitCode)
	}

	col := parseJSON(t, stdout)
	id, ok := col["id"].(string)
	if !ok || id == "" {
		t.Fatal("expected non-empty column id")
	}
	if col["key_name"] != "test_column" {
		t.Errorf("expected key_name 'test_column', got %v", col["key_name"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "get-column", "--dataset", dataset, "--id", id)
	if exitCode != 0 {
		t.Fatalf("get-column failed with exit code %d", exitCode)
	}
	got := parseJSON(t, stdout)
	if got["id"] != id {
		t.Errorf("expected id %q, got %v", id, got["id"])
	}

	stdout, _, exitCode = runCLIWithKey(t,
		"update-column",
		"--dataset", dataset,
		"--id", id,
		"--description", "updated description",
		"--type", "integer",
	)
	if exitCode != 0 {
		t.Fatalf("update-column failed with exit code %d", exitCode)
	}
	updated := parseJSON(t, stdout)
	if updated["description"] != "updated description" {
		t.Errorf("expected description 'updated description', got %v", updated["description"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "columns", "--dataset", dataset)
	if exitCode != 0 {
		t.Fatalf("columns failed with exit code %d", exitCode)
	}
	arr := parseJSONArray(t, stdout)
	found := false
	for _, item := range arr {
		c := item.(map[string]any)
		if c["id"] == id {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find column %q in list", id)
	}

	_, _, exitCode = runCLIWithKey(t, "delete-column", "--dataset", dataset, "--id", id)
	if exitCode != 0 {
		t.Fatalf("delete-column failed with exit code %d", exitCode)
	}
}
