package timefmt_test

import (
	"testing"
	"time"

	"github.com/LarsEckart/hccli/timefmt"
)

func TestParseTimeRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		// Plain integers
		{"plain 3600", "3600", 3600, false},
		{"plain 0", "0", 0, false},
		{"plain 86400", "86400", 86400, false},

		// Duration strings
		{"4 hours", "4 hours", 14400, false},
		{"1 hour", "1 hour", 3600, false},
		{"30 minutes", "30 minutes", 1800, false},
		{"1 minute", "1 minute", 60, false},
		{"1 day", "1 day", 86400, false},
		{"7 days", "7 days", 604800, false},
		{"1 week", "1 week", 604800, false},
		{"2 weeks", "2 weeks", 1209600, false},
		{"10 seconds", "10 seconds", 10, false},
		{"1 second", "1 second", 1, false},

		// Short unit forms
		{"30 min", "30 min", 1800, false},
		{"5 mins", "5 mins", 300, false},
		{"2 hr", "2 hr", 7200, false},
		{"3 hrs", "3 hrs", 10800, false},

		// Keywords
		{"last hour", "last hour", 3600, false},
		{"last day", "last day", 86400, false},
		{"last week", "last week", 604800, false},

		// Whitespace handling
		{"whitespace plain", "  3600  ", 3600, false},
		{"whitespace duration", " 4 hours ", 14400, false},

		// Error cases
		{"empty", "", 0, true},
		{"abc", "abc", 0, true},
		{"reversed", "hours 4", 0, true},
		{"negative", "-1 hours", 0, true},
		{"unknown unit", "4 fortnights", 0, true},
		{"decimal", "4.5 hours", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timefmt.ParseTimeRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeRange(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseTimeRange(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	nyLoc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("failed to load America/New_York location: %v", err)
	}

	// Precompute expected timestamps using today's date: 2026-02-12
	iso8601Want := time.Date(2026, 2, 12, 14, 30, 0, 0, time.UTC).Unix()
	dateTimeWant := time.Date(2026, 2, 12, 14, 30, 0, 0, time.UTC).Unix()
	dateTimeSecWant := time.Date(2026, 2, 12, 14, 30, 45, 0, time.UTC).Unix()
	dateOnlyWant := time.Date(2026, 2, 12, 0, 0, 0, 0, time.UTC).Unix()
	dateTimeNYWant := time.Date(2026, 2, 12, 14, 30, 0, 0, nyLoc).Unix()

	tests := []struct {
		name    string
		input   string
		loc     *time.Location
		want    int64
		wantErr bool
	}{
		// Unix epoch
		{"unix epoch", "1771077000", time.UTC, 1771077000, false},

		// ISO 8601
		{"iso8601 Z", "2026-02-12T14:30:00Z", time.UTC, iso8601Want, false},

		// Date + time
		{"date time UTC", "2026-02-12 14:30", time.UTC, dateTimeWant, false},

		// Date + time + seconds
		{"date time seconds UTC", "2026-02-12 14:30:45", time.UTC, dateTimeSecWant, false},

		// Date only
		{"date only UTC", "2026-02-12", time.UTC, dateOnlyWant, false},

		// Non-UTC location
		{"date time NY", "2026-02-12 14:30", nyLoc, dateTimeNYWant, false},

		// Nil loc defaults to UTC
		{"nil loc", "2026-02-12 14:30", nil, dateTimeWant, false},

		// Error cases
		{"empty", "", time.UTC, 0, true},
		{"not a date", "not-a-date", time.UTC, 0, true},
		{"invalid month", "2026-13-01", time.UTC, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timefmt.ParseTimestamp(tt.input, tt.loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimestamp(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseTimestamp(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}

	// Verify NY produces different timestamp than UTC for the same wall clock time
	utcTS, _ := timefmt.ParseTimestamp("2026-02-12 14:30", time.UTC)
	nyTS, _ := timefmt.ParseTimestamp("2026-02-12 14:30", nyLoc)
	if utcTS == nyTS {
		t.Errorf("expected different timestamps for UTC vs America/New_York, both got %d", utcTS)
	}
}
