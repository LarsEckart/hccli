package main_test

import (
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "--version")
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d, stderr=%s stdout=%s", exitCode, stderr, stdout)
	}
	if !strings.Contains(stdout, "hccli version ") {
		t.Fatalf("expected version output, got:\n%s", stdout)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got:\n%s", stderr)
	}
}
