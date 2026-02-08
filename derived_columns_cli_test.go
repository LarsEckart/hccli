package main_test

import (
	"testing"
)

func TestDerivedColumnCRUDCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-derived-column",
		"--dataset", dataset,
		"--alias", "test_derived_col",
		"--expression", "INT(1)",
		"--description", "test derived column",
	)
	if exitCode != 0 {
		t.Fatalf("create-derived-column failed with exit code %d", exitCode)
	}

	col := parseJSON(t, stdout)
	id, ok := col["id"].(string)
	if !ok || id == "" {
		t.Fatal("expected non-empty derived column id")
	}
	if col["alias"] != "test_derived_col" {
		t.Errorf("expected alias 'test_derived_col', got %v", col["alias"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "get-derived-column", "--dataset", dataset, "--id", id)
	if exitCode != 0 {
		t.Fatalf("get-derived-column failed with exit code %d", exitCode)
	}
	got := parseJSON(t, stdout)
	if got["id"] != id {
		t.Errorf("expected id %q, got %v", id, got["id"])
	}

	stdout, _, exitCode = runCLIWithKey(t,
		"update-derived-column",
		"--dataset", dataset,
		"--id", id,
		"--alias", "test_derived_col_updated",
		"--expression", "INT(2)",
		"--description", "updated description",
	)
	if exitCode != 0 {
		t.Fatalf("update-derived-column failed with exit code %d", exitCode)
	}
	updated := parseJSON(t, stdout)
	if updated["alias"] != "test_derived_col_updated" {
		t.Errorf("expected alias 'test_derived_col_updated', got %v", updated["alias"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "derived-columns", "--dataset", dataset)
	if exitCode != 0 {
		t.Fatalf("derived-columns failed with exit code %d", exitCode)
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
		t.Errorf("expected to find derived column %q in list", id)
	}

	_, _, exitCode = runCLIWithKey(t, "delete-derived-column", "--dataset", dataset, "--id", id)
	if exitCode != 0 {
		t.Fatalf("delete-derived-column failed with exit code %d", exitCode)
	}
}
