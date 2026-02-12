// Package timefmt provides human-friendly time parsing for CLI flags.
package timefmt

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseTimeRange parses a time range string into seconds.
// Supports plain integers ("3600"), duration strings ("4 hours"),
// and keywords ("last hour", "last day", "last week").
func ParseTimeRange(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty time range string")
	}

	lower := strings.ToLower(s)

	// Keywords
	switch lower {
	case "last hour":
		return 3600, nil
	case "last day":
		return 86400, nil
	case "last week":
		return 604800, nil
	}

	// Plain integer
	if n, err := strconv.Atoi(lower); err == nil {
		return n, nil
	}

	// Duration with units: "4 hours", "30 min"
	fields := strings.Fields(lower)
	if len(fields) != 2 {
		return 0, fmt.Errorf("unrecognized time range format: %q", s)
	}

	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid number in time range %q: %w", s, err)
	}
	if n < 0 {
		return 0, fmt.Errorf("time range must be positive: %q", s)
	}

	unit := fields[1]
	var multiplier int
	switch unit {
	case "second", "seconds":
		multiplier = 1
	case "minute", "minutes", "min", "mins":
		multiplier = 60
	case "hour", "hours", "hr", "hrs":
		multiplier = 3600
	case "day", "days":
		multiplier = 86400
	case "week", "weeks":
		multiplier = 604800
	default:
		return 0, fmt.Errorf("unknown time unit %q in %q", unit, s)
	}

	return n * multiplier, nil
}

// ParseTimestamp parses an absolute timestamp string into a Unix epoch (seconds).
// Supports Unix epoch integers, ISO 8601, "YYYY-MM-DD HH:MM[:SS]", and date-only.
// If loc is nil, defaults to time.UTC.
func ParseTimestamp(s string, loc *time.Location) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty timestamp string")
	}

	if loc == nil {
		loc = time.UTC
	}

	// Unix epoch integer
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n, nil
	}

	// ISO 8601 with timezone
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.Unix(), nil
	}

	// Date + time with seconds
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, loc); err == nil {
		return t.Unix(), nil
	}

	// Date + time (no seconds)
	if t, err := time.ParseInLocation("2006-01-02 15:04", s, loc); err == nil {
		return t.Unix(), nil
	}

	// Date only
	if t, err := time.ParseInLocation("2006-01-02", s, loc); err == nil {
		return t.Unix(), nil
	}

	return 0, fmt.Errorf("unrecognized timestamp format: %q", s)
}
