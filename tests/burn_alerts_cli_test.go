package main_test

import (
	"fmt"
	"testing"
)

func TestBurnAlertCRUDCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	// Create a derived column to use as SLI
	stdout, _, exitCode := runCLIWithKey(t,
		"create-derived-column",
		"--dataset", dataset,
		"--alias", "test_ba_sli",
		"--expression", "BOOL(1)",
	)
	if exitCode != 0 {
		t.Fatalf("create-derived-column failed with exit code %d", exitCode)
	}
	dc := parseJSON(t, stdout)
	dcID := dc["id"].(string)
	defer func() {
		runCLIWithKey(t, "delete-derived-column", "--dataset", dataset, "--id", dcID)
	}()

	// Create SLO
	stdout, _, exitCode = runCLIWithKey(t,
		"create-slo",
		"--dataset", dataset,
		"--name", "test SLO for burn alerts",
		"--sli-alias", "test_ba_sli",
		"--time-period-days", "30",
		"--target-per-million", "999000",
	)
	if exitCode != 0 {
		t.Fatalf("create-slo failed with exit code %d", exitCode)
	}
	slo := parseJSON(t, stdout)
	sloID := slo["id"].(string)
	defer func() {
		runCLIWithKey(t, "delete-slo", "--dataset", dataset, "--id", sloID)
	}()

	recipientsJSON := `[{"type":"email","target":"hccli-test@example.com"}]`

	// Create burn alert (exhaustion_time)
	stdout, _, exitCode = runCLIWithKey(t,
		"create-burn-alert",
		"--dataset", dataset,
		"--slo-id", sloID,
		"--alert-type", "exhaustion_time",
		"--exhaustion-minutes", "120",
		"--description", "test burn alert",
		"--recipients-json", recipientsJSON,
	)
	if exitCode != 0 {
		t.Fatalf("create-burn-alert failed with exit code %d", exitCode)
	}
	ba := parseJSON(t, stdout)
	baID, ok := ba["id"].(string)
	if !ok || baID == "" {
		t.Fatal("expected non-empty burn alert id")
	}
	if ba["alert_type"] != "exhaustion_time" {
		t.Errorf("expected alert_type 'exhaustion_time', got %v", ba["alert_type"])
	}
	if ba["description"] != "test burn alert" {
		t.Errorf("expected description 'test burn alert', got %v", ba["description"])
	}
	defer func() {
		runCLIWithKey(t, "delete-burn-alert", "--dataset", dataset, "--id", baID)
	}()

	// Get burn alert
	stdout, _, exitCode = runCLIWithKey(t, "get-burn-alert", "--dataset", dataset, "--id", baID)
	if exitCode != 0 {
		t.Fatalf("get-burn-alert failed with exit code %d", exitCode)
	}
	got := parseJSON(t, stdout)
	if got["id"] != baID {
		t.Errorf("expected id %q, got %v", baID, got["id"])
	}

	// Update burn alert
	stdout, _, exitCode = runCLIWithKey(t,
		"update-burn-alert",
		"--dataset", dataset,
		"--id", baID,
		"--alert-type", "exhaustion_time",
		"--exhaustion-minutes", "60",
		"--description", "updated burn alert",
		"--recipients-json", recipientsJSON,
	)
	if exitCode != 0 {
		t.Fatalf("update-burn-alert failed with exit code %d", exitCode)
	}
	updated := parseJSON(t, stdout)
	if updated["description"] != "updated burn alert" {
		t.Errorf("expected description 'updated burn alert', got %v", updated["description"])
	}
	exhaustionMin := updated["exhaustion_minutes"].(float64)
	if exhaustionMin != 60 {
		t.Errorf("expected exhaustion_minutes 60, got %v", exhaustionMin)
	}

	// List burn alerts
	stdout, _, exitCode = runCLIWithKey(t,
		"burn-alerts",
		"--dataset", dataset,
		"--slo-id", sloID,
	)
	if exitCode != 0 {
		t.Fatalf("burn-alerts failed with exit code %d", exitCode)
	}
	arr := parseJSONArray(t, stdout)
	found := false
	for _, item := range arr {
		a := item.(map[string]any)
		if a["id"] == baID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find burn alert %q in list", baID)
	}

	// Delete burn alert (explicit, before defer)
	_, stderr, exitCode := runCLIWithKey(t, "delete-burn-alert", "--dataset", dataset, "--id", baID)
	if exitCode != 0 {
		t.Fatalf("delete-burn-alert failed with exit code %d: %s", exitCode, stderr)
	}

	// Verify deleted — list should no longer contain it
	stdout, _, exitCode = runCLIWithKey(t, "burn-alerts", "--dataset", dataset, "--slo-id", sloID)
	if exitCode != 0 {
		t.Fatalf("burn-alerts (after delete) failed with exit code %d", exitCode)
	}
	arr = parseJSONArray(t, stdout)
	for _, item := range arr {
		a := item.(map[string]any)
		if a["id"] == baID {
			t.Errorf("burn alert %q should have been deleted but still in list", baID)
		}
	}
}

func TestBurnAlertBudgetRateCLI_Smoke(t *testing.T) {
	dataset := requireDataset(t)

	// Create a derived column to use as SLI
	stdout, _, exitCode := runCLIWithKey(t,
		"create-derived-column",
		"--dataset", dataset,
		"--alias", "test_ba_br_sli",
		"--expression", "BOOL(1)",
	)
	if exitCode != 0 {
		t.Fatalf("create-derived-column failed with exit code %d", exitCode)
	}
	dc := parseJSON(t, stdout)
	dcID := dc["id"].(string)
	defer func() {
		runCLIWithKey(t, "delete-derived-column", "--dataset", dataset, "--id", dcID)
	}()

	// Create SLO
	stdout, _, exitCode = runCLIWithKey(t,
		"create-slo",
		"--dataset", dataset,
		"--name", "test SLO for budget rate",
		"--sli-alias", "test_ba_br_sli",
		"--time-period-days", "30",
		"--target-per-million", "999000",
	)
	if exitCode != 0 {
		t.Fatalf("create-slo failed with exit code %d", exitCode)
	}
	slo := parseJSON(t, stdout)
	sloID := slo["id"].(string)
	defer func() {
		runCLIWithKey(t, "delete-slo", "--dataset", dataset, "--id", sloID)
	}()

	recipientsJSON := `[{"type":"email","target":"hccli-test@example.com"}]`

	// Create budget_rate burn alert
	stdout, _, exitCode = runCLIWithKey(t,
		"create-burn-alert",
		"--dataset", dataset,
		"--slo-id", sloID,
		"--alert-type", "budget_rate",
		"--budget-rate-window-minutes", "60",
		"--budget-rate-decrease-per-million", "10000",
		"--description", "budget rate test",
		"--recipients-json", recipientsJSON,
	)
	if exitCode != 0 {
		t.Fatalf("create-burn-alert (budget_rate) failed with exit code %d", exitCode)
	}
	ba := parseJSON(t, stdout)
	baID, ok := ba["id"].(string)
	if !ok || baID == "" {
		t.Fatal("expected non-empty burn alert id")
	}
	defer func() {
		runCLIWithKey(t, "delete-burn-alert", "--dataset", dataset, "--id", baID)
	}()

	if ba["alert_type"] != "budget_rate" {
		t.Errorf("expected alert_type 'budget_rate', got %v", ba["alert_type"])
	}
	windowMin := ba["budget_rate_window_minutes"].(float64)
	if windowMin != 60 {
		t.Errorf("expected budget_rate_window_minutes 60, got %v", windowMin)
	}
	threshold := ba["budget_rate_decrease_threshold_per_million"].(float64)
	if threshold != 10000 {
		t.Errorf("expected budget_rate_decrease_threshold_per_million 10000, got %v", threshold)
	}

	fmt.Println("budget_rate burn alert created and verified successfully")
}
