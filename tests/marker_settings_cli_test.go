package main_test

import (
	"testing"
)

func TestMarkerSettingCRUDCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-marker-setting",
		"--dataset", dataset,
		"--type", "test-setting",
		"--color", "#FF0000",
	)
	if exitCode != 0 {
		t.Fatalf("create-marker-setting failed with exit code %d", exitCode)
	}

	ms := parseJSON(t, stdout)
	id, ok := ms["id"].(string)
	if !ok || id == "" {
		t.Fatal("expected non-empty marker setting id")
	}
	if ms["type"] != "test-setting" {
		t.Errorf("expected type 'test-setting', got %v", ms["type"])
	}
	if ms["color"] != "#FF0000" {
		t.Errorf("expected color '#FF0000', got %v", ms["color"])
	}

	stdout, _, exitCode = runCLIWithKey(t,
		"update-marker-setting",
		"--dataset", dataset,
		"--id", id,
		"--type", "test-setting",
		"--color", "#00FF00",
	)
	if exitCode != 0 {
		t.Fatalf("update-marker-setting failed with exit code %d", exitCode)
	}
	updated := parseJSON(t, stdout)
	if updated["color"] != "#00FF00" {
		t.Errorf("expected color '#00FF00', got %v", updated["color"])
	}

	stdout, _, exitCode = runCLIWithKey(t, "marker-settings", "--dataset", dataset)
	if exitCode != 0 {
		t.Fatalf("marker-settings failed with exit code %d", exitCode)
	}
	arr := parseJSONArray(t, stdout)
	found := false
	for _, item := range arr {
		s := item.(map[string]any)
		if s["id"] == id {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find marker setting %q in list", id)
	}

	_, _, exitCode = runCLIWithKey(t, "delete-marker-setting", "--dataset", dataset, "--id", id)
	if exitCode != 0 {
		t.Fatalf("delete-marker-setting failed with exit code %d", exitCode)
	}
}
