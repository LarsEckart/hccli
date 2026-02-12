package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEmptyQueryResultShowsWarning(t *testing.T) {
	response := `{"id":"result-1","complete":true,"query_id":"q-1","data":{"series":[],"results":[]}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))
	defer srv.Close()

	stdout, stderr, code := runCLI(t,
		"--api-key", "fake-key",
		"--api-url", srv.URL,
		"create-query-result", "--dataset", "test-dataset", "--query-id", "q-1",
	)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}

	if !strings.Contains(stderr, "Query returned 0 results") {
		t.Errorf("expected empty results warning on stderr, got: %s", stderr)
	}
	if !strings.Contains(stderr, "Possible reasons") {
		t.Errorf("expected possible reasons on stderr, got: %s", stderr)
	}
	if !strings.Contains(stderr, "hccli columns --dataset test-dataset") {
		t.Errorf("expected dataset-specific column hint on stderr, got: %s", stderr)
	}
	if !strings.Contains(stderr, "breakdowns") {
		t.Errorf("expected breakdown hint on stderr, got: %s", stderr)
	}

	// stdout should still contain valid JSON
	result := parseJSON(t, stdout)
	if result["id"] != "result-1" {
		t.Errorf("expected result id 'result-1', got %v", result["id"])
	}
}

func TestNonEmptyQueryResultNoWarning(t *testing.T) {
	response := `{"id":"result-2","complete":true,"query_id":"q-2","data":{"series":[],"results":[{"data":{"count":42}}]}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))
	defer srv.Close()

	stdout, stderr, code := runCLI(t,
		"--api-key", "fake-key",
		"--api-url", srv.URL,
		"create-query-result", "--dataset", "test-dataset", "--query-id", "q-2",
	)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}

	if strings.Contains(stderr, "Query returned 0 results") {
		t.Errorf("expected no warning for non-empty results, got stderr: %s", stderr)
	}

	result := parseJSON(t, stdout)
	if result["id"] != "result-2" {
		t.Errorf("expected result id 'result-2', got %v", result["id"])
	}
}
