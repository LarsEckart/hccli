package cmd

import "testing"

func TestTraceUIURL_EscapesURLComponents(t *testing.T) {
	got := traceUIURL(
		"team/slug",
		"prod?env",
		"front/end?name=a b#frag",
		"trace/abc?x=1&foo=bar#frag space",
	)
	want := "https://ui.honeycomb.io/team%2Fslug/environments/prod%3Fenv/datasets/front%2Fend%3Fname=a%20b%23frag/trace?trace_id=trace%2Fabc%3Fx%3D1%26foo%3Dbar%23frag+space"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
