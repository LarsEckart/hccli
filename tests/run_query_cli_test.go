package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunQueryCLI(t *testing.T) {
	var createdQuery map[string]any
	var resultRequest map[string]any
	resultPolls := 0

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/1/queries/test-dataset":
			if err := json.NewDecoder(r.Body).Decode(&createdQuery); err != nil {
				t.Fatalf("failed to decode create-query body: %v", err)
			}
			response := map[string]any{"id": "query-1"}
			for k, v := range createdQuery {
				response[k] = v
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Fatalf("failed to encode create-query response: %v", err)
			}
		case r.Method == http.MethodPost && r.URL.Path == "/1/query_results/test-dataset":
			if err := json.NewDecoder(r.Body).Decode(&resultRequest); err != nil {
				t.Fatalf("failed to decode create-query-result body: %v", err)
			}
			if err := json.NewEncoder(w).Encode(map[string]any{
				"id":       "result-1",
				"complete": false,
				"query_id": "query-1",
			}); err != nil {
				t.Fatalf("failed to encode incomplete query result response: %v", err)
			}
		case r.Method == http.MethodGet && r.URL.Path == "/1/query_results/test-dataset/result-1":
			resultPolls++
			if err := json.NewEncoder(w).Encode(map[string]any{
				"id":       "result-1",
				"complete": true,
				"query_id": "query-1",
				"links": map[string]any{
					"query_url": "https://ui.honeycomb.io/team/datasets/test-dataset/result/result-1",
				},
				"data": map[string]any{
					"series": []any{},
					"results": []any{
						map[string]any{"data": map[string]any{"MAX(duration_ms)": 1234}},
					},
				},
			}); err != nil {
				t.Fatalf("failed to encode query result response: %v", err)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	stdout, stderr, code := runCLI(t,
		"--api-key", "fake-key",
		"--api-url", srv.URL,
		"run-query",
		"--dataset", "test-dataset",
		"--calculation-op", "MAX",
		"--calculation-column", "duration_ms",
		"--breakdown", "trace.trace_id",
		"--order", "MAX(duration_ms) desc",
		"--limit", "1",
		"--time-range", "30 minutes",
		"--poll-interval", "1",
		"--result-timeout", "5",
	)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr for non-empty results, got: %s", stderr)
	}

	if createdQuery["time_range"] != float64(1800) {
		t.Errorf("expected time_range 1800, got %v", createdQuery["time_range"])
	}
	if createdQuery["limit"] != float64(1) {
		t.Errorf("expected limit 1, got %v", createdQuery["limit"])
	}

	if resultRequest["query_id"] != "query-1" {
		t.Fatalf("expected query result request for query-1, got %v", resultRequest)
	}
	if resultPolls != 1 {
		t.Fatalf("expected one result poll, got %d", resultPolls)
	}

	result := parseJSON(t, stdout)
	if result["id"] != "result-1" {
		t.Errorf("expected result-1 output, got %v", result["id"])
	}
	if result["query_id"] != "query-1" {
		t.Errorf("expected query_id query-1 output, got %v", result["query_id"])
	}
	if result["complete"] != true {
		t.Errorf("expected complete output, got %v", result["complete"])
	}
}

func TestRunQueryInvalidArgumentsFailBeforeAPI(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		http.Error(w, "should not be called", http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, stderr, code := runCLI(t,
		"--api-key", "fake-key",
		"--api-url", srv.URL,
		"run-query",
		"--dataset", "test-dataset",
		"--calculation-op", "COUNT",
		"--breakdown", "service.name",
		"--order", "duration_ms desc",
	)
	if code == 0 {
		t.Fatal("expected non-zero exit code for invalid run-query arguments")
	}
	if called {
		t.Fatal("expected invalid run-query arguments to fail before calling API")
	}
	if !strings.Contains(stderr, "invalid order") {
		t.Errorf("expected invalid order error, got: %s", stderr)
	}
}

func TestRunQueryFromQueryJSONStdinCLI(t *testing.T) {
	var createdQuery map[string]any
	var resultRequest map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/1/queries/test-dataset":
			if err := json.NewDecoder(r.Body).Decode(&createdQuery); err != nil {
				t.Fatalf("failed to decode create-query body: %v", err)
			}
			if err := json.NewEncoder(w).Encode(map[string]any{"id": "query-json-1"}); err != nil {
				t.Fatalf("failed to encode create-query response: %v", err)
			}
		case r.Method == http.MethodPost && r.URL.Path == "/1/query_results/test-dataset":
			if err := json.NewDecoder(r.Body).Decode(&resultRequest); err != nil {
				t.Fatalf("failed to decode create-query-result body: %v", err)
			}
			if err := json.NewEncoder(w).Encode(map[string]any{
				"id":       "result-json-1",
				"complete": true,
				"query_id": "query-json-1",
				"data": map[string]any{
					"results": []any{map[string]any{"data": map[string]any{"COUNT": 10}}},
				},
			}); err != nil {
				t.Fatalf("failed to encode query result response: %v", err)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	stdout, stderr, code := runCLIWithOptions(t,
		cliRunOptions{Stdin: `{"calculations":[{"op":"COUNT"}],"future_field":{"preserved":true}}`},
		"--api-key", "fake-key",
		"--api-url", srv.URL,
		"run-query",
		"--dataset", "test-dataset",
		"--query-json", "-",
	)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", code, stdout, stderr)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr for non-empty results, got: %s", stderr)
	}

	futureField, ok := createdQuery["future_field"].(map[string]any)
	if !ok || futureField["preserved"] != true {
		t.Fatalf("expected future field to be preserved, got %v", createdQuery["future_field"])
	}
	if resultRequest["query_id"] != "query-json-1" {
		t.Fatalf("expected query result request for query-json-1, got %v", resultRequest)
	}

	result := parseJSON(t, stdout)
	if result["id"] != "result-json-1" {
		t.Errorf("expected result-json-1 output, got %v", result["id"])
	}
}
