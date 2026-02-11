package main_test

import (
	"testing"
)

func TestSLOCRUDCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	// Create a derived column to use as SLI
	stdout, _, exitCode := runCLIWithKey(t,
		"create-derived-column",
		"--dataset", dataset,
		"--alias", "test_slo_sli",
		"--expression", "BOOL(1)",
	)
	if exitCode != 0 {
		t.Fatalf("create-derived-column failed with exit code %d", exitCode)
	}
	dc := parseJSON(t, stdout)
	dcID, ok := dc["id"].(string)
	if !ok || dcID == "" {
		t.Fatal("expected non-empty derived column id")
	}
	defer func() {
		runCLIWithKey(t, "delete-derived-column", "--dataset", dataset, "--id", dcID)
	}()

	// Create SLO
	stdout, _, exitCode = runCLIWithKey(t,
		"create-slo",
		"--dataset", dataset,
		"--name", "test SLO",
		"--description", "a test SLO",
		"--sli-alias", "test_slo_sli",
		"--time-period-days", "30",
		"--target-per-million", "999000",
	)
	if exitCode != 0 {
		t.Fatalf("create-slo failed with exit code %d", exitCode)
	}
	slo := parseJSON(t, stdout)
	sloID, ok := slo["id"].(string)
	if !ok || sloID == "" {
		t.Fatal("expected non-empty SLO id")
	}
	if slo["name"] != "test SLO" {
		t.Errorf("expected name 'test SLO', got %v", slo["name"])
	}
	defer func() {
		runCLIWithKey(t, "delete-slo", "--dataset", dataset, "--id", sloID)
	}()

	// Get SLO
	stdout, _, exitCode = runCLIWithKey(t, "get-slo", "--dataset", dataset, "--id", sloID)
	if exitCode != 0 {
		t.Fatalf("get-slo failed with exit code %d", exitCode)
	}
	got := parseJSON(t, stdout)
	if got["id"] != sloID {
		t.Errorf("expected id %q, got %v", sloID, got["id"])
	}

	// Update SLO
	stdout, _, exitCode = runCLIWithKey(t,
		"update-slo",
		"--dataset", dataset,
		"--id", sloID,
		"--name", "updated SLO",
		"--description", "updated description",
		"--sli-alias", "test_slo_sli",
		"--time-period-days", "30",
		"--target-per-million", "990000",
	)
	if exitCode != 0 {
		t.Fatalf("update-slo failed with exit code %d", exitCode)
	}
	updated := parseJSON(t, stdout)
	if updated["name"] != "updated SLO" {
		t.Errorf("expected name 'updated SLO', got %v", updated["name"])
	}

	// List SLOs
	stdout, _, exitCode = runCLIWithKey(t, "slos", "--dataset", dataset)
	if exitCode != 0 {
		t.Fatalf("slos failed with exit code %d", exitCode)
	}
	arr := parseJSONArray(t, stdout)
	found := false
	for _, item := range arr {
		s := item.(map[string]any)
		if s["id"] == sloID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find SLO %q in list", sloID)
	}

	// Delete SLO (explicit, before defer)
	_, _, exitCode = runCLIWithKey(t, "delete-slo", "--dataset", dataset, "--id", sloID)
	if exitCode != 0 {
		t.Fatalf("delete-slo failed with exit code %d", exitCode)
	}
}
