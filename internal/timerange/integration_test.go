package timerange_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/timerange"
)

func TestRoundTrip_RFC3339(t *testing.T) {
	r, err := timerange.Parse("2024-01-01T00:00:00Z", "2024-12-31T23:59:59Z")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	cases := []struct {
		ts   string
		want bool
	}{
		{"2024-06-01T12:00:00Z", true},
		{"2023-12-31T23:59:59Z", false},
		{"2025-01-01T00:00:00Z", false},
		{"2024-01-01T00:00:00Z", true},
		{"2024-12-31T23:59:59Z", true},
	}

	for _, c := range cases {
		entry := map[string]any{"timestamp": c.ts}
		got := r.Match(entry, "timestamp")
		if got != c.want {
			t.Errorf("Match(%q) = %v, want %v", c.ts, got, c.want)
		}
	}
}

func TestRoundTrip_DateOnly(t *testing.T) {
	r, err := timerange.Parse("2024-03-01", "2024-03-31")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	entry := map[string]any{"ts": "2024-03-15"}
	if !r.Match(entry, "ts") {
		t.Error("expected match for date within range")
	}

	outside := map[string]any{"ts": "2024-04-01"}
	if r.Match(outside, "ts") {
		t.Error("expected no match for date outside range")
	}
}
