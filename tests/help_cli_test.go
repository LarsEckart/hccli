package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHelpOutputConventions(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains []string
		excludes []string
	}{
		{
			name: "required flags and examples are visible",
			args: []string{"get-board", "--help"},
			contains: []string{
				"Examples:",
				"hccli get-board --id brd_123",
				"Board ID (required)",
			},
		},
		{
			name: "query polling timeout does not collide with global HTTP timeout",
			args: []string{"run-query", "--help"},
			contains: []string{
				"--result-timeout int",
				"HTTP request timeout in seconds",
			},
			excludes: []string{"Maximum seconds to wait for results"},
		},
		{
			name: "auth login hides duplicated global auth flags",
			args: []string{"auth", "login", "--help"},
			contains: []string{
				"Examples:",
				"--profile string, -p string",
				"Global --api-url and --timeout may be passed before auth login",
			},
			excludes: []string{"GLOBAL OPTIONS:"},
		},
		{
			name: "optional zero-valued ints do not show misleading defaults",
			args: []string{"update-marker", "--help"},
			contains: []string{
				"Examples:",
				"--start-time int",
			},
			excludes: []string{"default: 0"},
		},
		{
			name: "delete help states safety behavior",
			args: []string{"delete-board", "--help"},
			contains: []string{
				"Deletes immediately; no confirmation prompt.",
				"Examples:",
				"Board ID (required)",
			},
		},
		{
			name: "help command supports help flag",
			args: []string{"help", "--help"},
			contains: []string{
				"Examples:",
				"hccli auth status",
			},
			excludes: []string{"no help topic"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, code := runCLI(t, tc.args...)
			if code != 0 {
				t.Fatalf("help command failed with code %d\nstderr: %s", code, stderr)
			}
			for _, want := range tc.contains {
				if !strings.Contains(stdout, want) {
					t.Fatalf("expected help to contain %q\noutput: %s", want, stdout)
				}
			}
			for _, unwanted := range tc.excludes {
				if strings.Contains(stdout, unwanted) {
					t.Fatalf("expected help not to contain %q\noutput: %s", unwanted, stdout)
				}
			}
		})
	}
}

func TestDeleteCommandPrintsJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/1/boards/brd_123" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	stdout, stderr, code := runCLI(t,
		"--api-key", "fake-key",
		"--api-url", server.URL,
		"delete-board",
		"--id", "brd_123",
	)
	if code != 0 {
		t.Fatalf("delete-board failed with code %d\nstderr: %s", code, stderr)
	}

	var out map[string]any
	if err := json.Unmarshal([]byte(stdout), &out); err != nil {
		t.Fatalf("delete output is not JSON: %v\noutput: %s", err, stdout)
	}
	if out["deleted"] != true || out["resource"] != "board" || out["id"] != "brd_123" {
		t.Fatalf("unexpected delete output: %s", stdout)
	}
}
