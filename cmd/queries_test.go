package cmd

import "testing"

func TestParseFilterAllowsRepeatedWhitespace(t *testing.T) {
	filter, err := parseFilter("duration_ms   > 100")
	if err != nil {
		t.Fatalf("parseFilter returned error: %v", err)
	}

	if filter.Column != "duration_ms" {
		t.Fatalf("expected column %q, got %q", "duration_ms", filter.Column)
	}
	if filter.Op != ">" {
		t.Fatalf("expected op %q, got %q", ">", filter.Op)
	}
	if filter.Value != int64(100) {
		t.Fatalf("expected value %v (%T), got %v (%T)", int64(100), int64(100), filter.Value, filter.Value)
	}
}

func TestParseFilterPreservesValueTextAfterOperator(t *testing.T) {
	filter, err := parseFilter("message   contains   hello  world")
	if err != nil {
		t.Fatalf("parseFilter returned error: %v", err)
	}

	if filter.Column != "message" {
		t.Fatalf("expected column %q, got %q", "message", filter.Column)
	}
	if filter.Op != "contains" {
		t.Fatalf("expected op %q, got %q", "contains", filter.Op)
	}
	if filter.Value != "hello  world" {
		t.Fatalf("expected value %q, got %q", "hello  world", filter.Value)
	}
}
