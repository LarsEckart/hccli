package api_test

import (
	"testing"
	"time"

	"github.com/LarsEckart/hccli/api"
)

func TestNewClientUsesDefaultTimeoutWhenZero(t *testing.T) {
	client := api.NewClient("key", 0)

	if client.HTTP.Timeout != api.DefaultTimeout {
		t.Fatalf("expected timeout %s, got %s", api.DefaultTimeout, client.HTTP.Timeout)
	}
}

func TestNewClientUsesExplicitTimeout(t *testing.T) {
	timeout := 5 * time.Second
	client := api.NewClient("key", timeout)

	if client.HTTP.Timeout != timeout {
		t.Fatalf("expected timeout %s, got %s", timeout, client.HTTP.Timeout)
	}
}
