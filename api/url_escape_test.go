package api_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/LarsEckart/hccli/api"
)

type captureTransport struct {
	requestURI string
	rawQuery   string
	body       string
}

func (t *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.requestURI = req.URL.RequestURI()
	t.rawQuery = req.URL.RawQuery
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(t.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func TestGetColumnEscapesResourcePathSegments(t *testing.T) {
	transport := &captureTransport{body: `{}`}
	client := api.NewClient("test-key", time.Second)
	client.BaseURL = "https://example.test"
	client.HTTP = &http.Client{Transport: transport}

	_, err := client.GetColumn(t.Context(), "dataset/with?reserved#chars", "column/with?reserved")
	if err != nil {
		t.Fatalf("get column: %v", err)
	}

	want := "/1/columns/dataset%2Fwith%3Freserved%23chars/column%2Fwith%3Freserved"
	if transport.requestURI != want {
		t.Fatalf("request URI = %q, want %q", transport.requestURI, want)
	}
}

func TestListBurnAlertsEscapesPathAndQueryValues(t *testing.T) {
	transport := &captureTransport{body: `[]`}
	client := api.NewClient("test-key", time.Second)
	client.BaseURL = "https://example.test"
	client.HTTP = &http.Client{Transport: transport}

	_, err := client.ListBurnAlerts(t.Context(), "dataset/with?reserved#chars", "abc&foo=bar")
	if err != nil {
		t.Fatalf("list burn alerts: %v", err)
	}

	wantURI := "/1/burn_alerts/dataset%2Fwith%3Freserved%23chars?slo_id=abc%26foo%3Dbar"
	if transport.requestURI != wantURI {
		t.Fatalf("request URI = %q, want %q", transport.requestURI, wantURI)
	}

	wantQuery := "slo_id=abc%26foo%3Dbar"
	if transport.rawQuery != wantQuery {
		t.Fatalf("raw query = %q, want %q", transport.rawQuery, wantQuery)
	}
}
