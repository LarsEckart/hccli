package main_test

import (
	"testing"
)

func createTestBoard(t *testing.T) string {
	t.Helper()
	stdout, _, exitCode := runCLIWithKey(t, "create-board", "--name", "hccli view test board")
	if exitCode != 0 {
		t.Fatalf("create-board failed with exit code %d", exitCode)
	}
	board := parseJSON(t, stdout)
	return board["id"].(string)
}

func deleteTestBoard(t *testing.T, id string) {
	t.Helper()
	runCLIWithKey(t, "delete-board", "--id", id)
}

func TestCreateBoardViewCLI_Smoke(t *testing.T) {
	boardID := createTestBoard(t)
	defer deleteTestBoard(t, boardID)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-board-view",
		"--board-id", boardID,
		"--name", "test view",
		"--filter-column", "status",
		"--filter-op", "=",
		"--filter-value", "200",
	)
	if exitCode != 0 {
		t.Fatalf("create-board-view failed with exit code %d", exitCode)
	}

	view := parseJSON(t, stdout)
	if view["id"] == nil || view["id"] == "" {
		t.Fatal("expected non-empty view id")
	}
	if view["name"] != "test view" {
		t.Errorf("expected name 'test view', got %v", view["name"])
	}
}

func TestListBoardViewsCLI_Smoke(t *testing.T) {
	boardID := createTestBoard(t)
	defer deleteTestBoard(t, boardID)

	runCLIWithKey(t,
		"create-board-view",
		"--board-id", boardID,
		"--name", "list test view",
		"--filter-column", "status",
		"--filter-op", "=",
		"--filter-value", "200",
	)

	stdout, _, exitCode := runCLIWithKey(t, "board-views", "--board-id", boardID)
	if exitCode != 0 {
		t.Fatalf("board-views failed with exit code %d", exitCode)
	}

	arr := parseJSONArray(t, stdout)
	if len(arr) < 1 {
		t.Fatalf("expected at least 1 view, got %d", len(arr))
	}
}

func TestGetBoardViewCLI_Smoke(t *testing.T) {
	boardID := createTestBoard(t)
	defer deleteTestBoard(t, boardID)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-board-view",
		"--board-id", boardID,
		"--name", "get test view",
		"--filter-column", "status",
		"--filter-op", "=",
		"--filter-value", "200",
	)
	if exitCode != 0 {
		t.Fatalf("create-board-view failed with exit code %d", exitCode)
	}
	view := parseJSON(t, stdout)
	viewID := view["id"].(string)

	stdout, _, exitCode = runCLIWithKey(t, "get-board-view", "--board-id", boardID, "--view-id", viewID)
	if exitCode != 0 {
		t.Fatalf("get-board-view failed with exit code %d", exitCode)
	}

	got := parseJSON(t, stdout)
	if got["id"] != viewID {
		t.Errorf("expected view id %q, got %v", viewID, got["id"])
	}
	if got["name"] != "get test view" {
		t.Errorf("expected name 'get test view', got %v", got["name"])
	}
}

func TestUpdateBoardViewCLI_Smoke(t *testing.T) {
	boardID := createTestBoard(t)
	defer deleteTestBoard(t, boardID)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-board-view",
		"--board-id", boardID,
		"--name", "original view",
		"--filter-column", "status",
		"--filter-op", "=",
		"--filter-value", "200",
	)
	if exitCode != 0 {
		t.Fatalf("create-board-view failed with exit code %d", exitCode)
	}
	view := parseJSON(t, stdout)
	viewID := view["id"].(string)

	stdout, _, exitCode = runCLIWithKey(t,
		"update-board-view",
		"--board-id", boardID,
		"--view-id", viewID,
		"--name", "updated view",
		"--filter-column", "status",
		"--filter-op", "!=",
		"--filter-value", "500",
	)
	if exitCode != 0 {
		t.Fatalf("update-board-view failed with exit code %d", exitCode)
	}

	updated := parseJSON(t, stdout)
	if updated["name"] != "updated view" {
		t.Errorf("expected name 'updated view', got %v", updated["name"])
	}
}

func TestDeleteBoardViewCLI_Smoke(t *testing.T) {
	boardID := createTestBoard(t)
	defer deleteTestBoard(t, boardID)

	stdout, _, exitCode := runCLIWithKey(t,
		"create-board-view",
		"--board-id", boardID,
		"--name", "view to delete",
		"--filter-column", "status",
		"--filter-op", "=",
		"--filter-value", "200",
	)
	if exitCode != 0 {
		t.Fatalf("create-board-view failed with exit code %d", exitCode)
	}
	view := parseJSON(t, stdout)
	viewID := view["id"].(string)

	_, _, exitCode = runCLIWithKey(t, "delete-board-view", "--board-id", boardID, "--view-id", viewID)
	if exitCode != 0 {
		t.Fatalf("delete-board-view failed with exit code %d", exitCode)
	}

	_, _, exitCode = runCLIWithKey(t, "get-board-view", "--board-id", boardID, "--view-id", viewID)
	if exitCode == 0 {
		t.Error("expected non-zero exit code when getting deleted view")
	}
}
