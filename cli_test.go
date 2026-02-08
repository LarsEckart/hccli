package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "hccli-test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	binaryPath = filepath.Join(dir, "hccli")
	build := exec.CommandContext(context.Background(), "go", "build", "-o", binaryPath, ".")
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build binary: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func runCLI(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), binaryPath, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run CLI: %v", err)
		}
	}
	return stdout.String(), stderr.String(), exitCode
}

func runCLIWithKey(t *testing.T, args ...string) (string, string, int) {
	t.Helper()
	apiKey := os.Getenv("HONEYCOMB_API_KEY")
	if apiKey == "" {
		t.Skip("HONEYCOMB_API_KEY not set, skipping smoke test")
	}
	fullArgs := append([]string{"--api-key", apiKey}, args...)
	return runCLI(t, fullArgs...)
}

func requireDataset(t *testing.T) string {
	t.Helper()
	ds := os.Getenv("HONEYCOMB_DATASET")
	if ds == "" {
		t.Skip("HONEYCOMB_DATASET not set, skipping smoke test")
	}
	return ds
}

func parseJSON(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("failed to parse JSON output: %v\noutput: %s", err, s)
	}
	return m
}

func parseJSONArray(t *testing.T, s string) []any {
	t.Helper()
	var arr []any
	if err := json.Unmarshal([]byte(s), &arr); err != nil {
		t.Fatalf("failed to parse JSON array output: %v\noutput: %s", err, s)
	}
	return arr
}

func TestMissingAPIKeyShowsError(t *testing.T) {
	cmd := exec.CommandContext(t.Context(), binaryPath, "auth")
	cmd.Env = filterEnv("HONEYCOMB_API_KEY")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit code when api-key is missing")
	}
	if !strings.Contains(string(out), "api-key") {
		t.Fatalf("expected error mentioning api-key, got: %s", out)
	}
}

func filterEnv(exclude string) []string {
	var env []string
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, exclude+"=") {
			env = append(env, e)
		}
	}
	return env
}
