package main_test

import (
	"testing"
)

func TestListBoardsCLI_Smoke(t *testing.T) {
	stdout, _, exitCode := runCLIWithKey(t, "boards")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	arr := parseJSONArray(t, stdout)
	if arr == nil {
		t.Error("expected non-nil boards array")
	}
}

func TestCreateAndDeleteBoardCLI_Smoke(t *testing.T) {
	stdout, _, exitCode := runCLIWithKey(t, "create-board", "--name", "hccli smoke test board")
	if exitCode != 0 {
		t.Fatalf("create-board failed with exit code %d", exitCode)
	}

	board := parseJSON(t, stdout)
	id, ok := board["id"].(string)
	if !ok || id == "" {
		t.Fatal("expected non-empty board id")
	}
	if board["name"] != "hccli smoke test board" {
		t.Errorf("expected name 'hccli smoke test board', got %v", board["name"])
	}

	_, _, exitCode = runCLIWithKey(t, "delete-board", "--id", id)
	if exitCode != 0 {
		t.Fatalf("delete-board failed with exit code %d", exitCode)
	}
}

func TestGetBoardCLI_Smoke(t *testing.T) {
	stdout, _, exitCode := runCLIWithKey(t, "create-board", "--name", "hccli get board test")
	if exitCode != 0 {
		t.Fatalf("create-board failed with exit code %d", exitCode)
	}
	board := parseJSON(t, stdout)
	id := board["id"].(string)
	defer runCLIWithKey(t, "delete-board", "--id", id)

	stdout, _, exitCode = runCLIWithKey(t, "get-board", "--id", id)
	if exitCode != 0 {
		t.Fatalf("get-board failed with exit code %d", exitCode)
	}

	got := parseJSON(t, stdout)
	if got["id"] != id {
		t.Errorf("expected board id %q, got %v", id, got["id"])
	}
}

func TestUpdateBoardCLI_Smoke(t *testing.T) {
	stdout, _, exitCode := runCLIWithKey(t, "create-board", "--name", "hccli update board test")
	if exitCode != 0 {
		t.Fatalf("create-board failed with exit code %d", exitCode)
	}
	board := parseJSON(t, stdout)
	id := board["id"].(string)
	defer runCLIWithKey(t, "delete-board", "--id", id)

	stdout, _, exitCode = runCLIWithKey(t, "update-board", "--id", id, "--name", "updated name", "--description", "updated desc")
	if exitCode != 0 {
		t.Fatalf("update-board failed with exit code %d", exitCode)
	}

	updated := parseJSON(t, stdout)
	if updated["name"] != "updated name" {
		t.Errorf("expected name 'updated name', got %v", updated["name"])
	}
	if updated["description"] != "updated desc" {
		t.Errorf("expected description 'updated desc', got %v", updated["description"])
	}
}

func TestUpdateBoardPanelsJsonCLI_Smoke(t *testing.T) {
	stdout, _, exitCode := runCLIWithKey(t, "create-board", "--name", "hccli panels-json test")
	if exitCode != 0 {
		t.Fatalf("create-board failed with exit code %d", exitCode)
	}
	board := parseJSON(t, stdout)
	id := board["id"].(string)
	defer runCLIWithKey(t, "delete-board", "--id", id)

	panelsJSON := `[{"type":"text","text_panel":{"content":"hello"}}]`
	stdout, _, exitCode = runCLIWithKey(t, "update-board", "--id", id, "--name", "hccli panels-json test", "--panels-json", panelsJSON)
	if exitCode != 0 {
		t.Fatalf("update-board with panels-json failed with exit code %d", exitCode)
	}

	updated := parseJSON(t, stdout)
	panels, ok := updated["panels"].([]any)
	if !ok || len(panels) != 1 {
		t.Fatalf("expected 1 panel, got %v", updated["panels"])
	}

	emptyPanelsJSON := `[]`
	stdout, _, exitCode = runCLIWithKey(t, "update-board", "--id", id, "--name", "hccli panels-json test", "--panels-json", emptyPanelsJSON)
	if exitCode != 0 {
		t.Fatalf("update-board with empty panels-json failed with exit code %d", exitCode)
	}

	updated = parseJSON(t, stdout)
	panels, ok = updated["panels"].([]any)
	if ok && len(panels) != 0 {
		t.Errorf("expected 0 panels after removal, got %d", len(panels))
	}
}
