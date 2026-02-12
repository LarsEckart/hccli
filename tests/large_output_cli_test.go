package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLargeOutputWritesTempFile(t *testing.T) {
	// 392 entries produce ~30.7KB after indentation — just over the 30KB threshold.
	large := buildColumnsJSON(392)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, large)
	}))
	defer srv.Close()

	stdout, stderr, code := runCLI(t,
		"--api-key", "fake-key",
		"--api-url", srv.URL,
		"columns", "--dataset", "test-dataset",
	)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\nstderr: %s", code, stderr)
	}

	// stderr should contain the warning and a temp file path
	if !strings.Contains(stderr, "Output is large") {
		t.Errorf("expected large output warning on stderr, got: %s", stderr)
	}
	if !strings.Contains(stderr, "Full output written to:") {
		t.Errorf("expected temp file path on stderr, got: %s", stderr)
	}
	if !strings.Contains(stderr, "To reduce output size") {
		t.Errorf("expected suggestions on stderr, got: %s", stderr)
	}

	// Extract temp file path and verify it exists with correct content
	for _, line := range strings.Split(stderr, "\n") {
		if strings.Contains(line, "Full output written to:") {
			path := strings.TrimSpace(strings.SplitN(line, "Full output written to:", 2)[1])
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read temp file %s: %v", path, err)
			}
			if len(data) != len(stdout) {
				t.Errorf("temp file size (%d) != stdout size (%d)", len(data), len(stdout))
			}
			os.Remove(path)
			break
		}
	}

	// stdout should still contain the full JSON
	if len(stdout) < 30*1024 {
		t.Errorf("expected stdout > 30KB, got %d bytes", len(stdout))
	}
}

func TestSmallOutputNoWarning(t *testing.T) {
	// 391 entries produce ~30.6KB — just under the 30KB threshold.
	small := buildColumnsJSON(391)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, small)
	}))
	defer srv.Close()

	_, stderr, code := runCLI(t,
		"--api-key", "fake-key",
		"--api-url", srv.URL,
		"columns", "--dataset", "test-dataset",
	)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\nstderr: %s", code, stderr)
	}

	if strings.Contains(stderr, "Output is large") {
		t.Errorf("expected no warning for small output, got stderr: %s", stderr)
	}
}

// buildColumnsJSON creates a JSON array of fake columns.
// 392 entries produce ~30.7KB after re-indentation (just over the 30KB threshold).
// 391 entries produce ~30.6KB (just under).
func buildColumnsJSON(n int) string {
	var b strings.Builder
	b.WriteString("[")
	for i := range n {
		if i > 0 {
			b.WriteString(",")
		}
		fmt.Fprintf(&b, `{"id":"col-%d","key_name":"name_%d","type":"string"}`, i, i)
	}
	b.WriteString("]")
	return b.String()
}
