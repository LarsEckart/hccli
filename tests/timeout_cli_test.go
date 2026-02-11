package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTimeoutFlag(t *testing.T) {
	// Server that stalls until closed.
	done := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-done
	}))
	defer func() {
		close(done)
		server.Close()
	}()

	start := time.Now()
	_, stderr, exitCode := runCLI(t,
		"--api-key", "test-key",
		"--api-url", server.URL,
		"--timeout", "1",
		"auth",
	)
	elapsed := time.Since(start)

	if exitCode == 0 {
		t.Fatal("expected non-zero exit code on timeout")
	}

	if elapsed > 5*time.Second {
		t.Fatalf("expected timeout after ~1s, but took %s", elapsed)
	}

	lower := strings.ToLower(stderr)
	if !strings.Contains(lower, "timeout") && !strings.Contains(lower, "deadline") {
		t.Fatalf("expected timeout/deadline error, got: %s", stderr)
	}
}
