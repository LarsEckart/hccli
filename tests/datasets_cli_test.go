package main_test

import (
	"fmt"
	"testing"
	"time"
)

func TestDatasetCRUDCLI_Smoke(t *testing.T) {
	if os.Getenv("RUN_DATASET_TESTS") != "1" {
		t.Skip("set RUN_DATASET_TESTS=1 to run dataset integration tests")
	}
	_ = requireDataset(t)

	name := fmt.Sprintf("hccli-test-%d", time.Now().UnixNano())

	stdout, _, exitCode := runCLIWithKey(t,
		"create-dataset",
		"--name", name,
		"--description", "test dataset",
		"--expand-json-depth", "2",
	)
	if exitCode != 0 {
		t.Fatalf("create-dataset failed with exit code %d", exitCode)
	}

	ds := parseJSON(t, stdout)
	slug, ok := ds["slug"].(string)
	if !ok || slug == "" {
		t.Fatal("expected non-empty dataset slug")
	}
	if ds["name"] != name {
		t.Errorf("expected name %q, got %v", name, ds["name"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "get-dataset", "--slug", slug)
	if exitCode != 0 {
		t.Fatalf("get-dataset failed with exit code %d", exitCode)
	}
	got := parseJSON(t, stdout)
	if got["slug"] != slug {
		t.Errorf("expected slug %q, got %v", slug, got["slug"])
	}

	stdout, _, exitCode = runCLIWithKey(t,
		"update-dataset",
		"--slug", slug,
		"--description", "updated description",
		"--expand-json-depth", "3",
	)
	if exitCode != 0 {
		t.Fatalf("update-dataset failed with exit code %d", exitCode)
	}
	updated := parseJSON(t, stdout)
	if updated["description"] != "updated description" {
		t.Errorf("expected description 'updated description', got %v", updated["description"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "datasets")
	if exitCode != 0 {
		t.Fatalf("datasets failed with exit code %d", exitCode)
	}
	arr := parseJSONArray(t, stdout)
	found := false
	for _, item := range arr {
		d := item.(map[string]any)
		if d["slug"] == slug {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find dataset %q in list", slug)
	}

	_, _, exitCode = runCLIWithKey(t,
		"update-dataset",
		"--slug", slug,
		"--description", "updated description",
		"--expand-json-depth", "3",
		"--delete-protected=false",
	)
	if exitCode != 0 {
		t.Fatalf("disable delete protection failed with exit code %d", exitCode)
	}

	_, _, exitCode = runCLIWithKey(t, "delete-dataset", "--slug", slug)
	if exitCode != 0 {
		t.Fatalf("delete-dataset failed with exit code %d", exitCode)
	}
}
