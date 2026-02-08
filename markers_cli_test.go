package main_test

import (
	"testing"
)

func TestMarkerCRUDCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-marker",
		"--dataset", dataset,
		"--message", "test deploy #1",
		"--type", "deploy",
	)
	if exitCode != 0 {
		t.Fatalf("create-marker failed with exit code %d", exitCode)
	}

	marker := parseJSON(t, stdout)
	id, ok := marker["id"].(string)
	if !ok || id == "" {
		t.Fatal("expected non-empty marker id")
	}
	if marker["message"] != "test deploy #1" {
		t.Errorf("expected message 'test deploy #1', got %v", marker["message"])
	}
	if marker["type"] != "deploy" {
		t.Errorf("expected type 'deploy', got %v", marker["type"])
	}

	stdout, _, exitCode = runCLIWithKey(t,
		"update-marker",
		"--dataset", dataset,
		"--id", id,
		"--message", "updated deploy #1",
		"--type", "deploy",
	)
	if exitCode != 0 {
		t.Fatalf("update-marker failed with exit code %d", exitCode)
	}
	updated := parseJSON(t, stdout)
	if updated["message"] != "updated deploy #1" {
		t.Errorf("expected message 'updated deploy #1', got %v", updated["message"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "markers", "--dataset", dataset)
	if exitCode != 0 {
		t.Fatalf("markers failed with exit code %d", exitCode)
	}
	arr := parseJSONArray(t, stdout)
	found := false
	for _, item := range arr {
		m := item.(map[string]any)
		if m["id"] == id {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find marker %q in list", id)
	}

	_, _, exitCode = runCLIWithKey(t, "delete-marker", "--dataset", dataset, "--id", id)
	if exitCode != 0 {
		t.Fatalf("delete-marker failed with exit code %d", exitCode)
	}
}
