package timerange

import (
	"testing"
	"time"
)

func mustTime(s string) time.Time {
	t, err := parseTime(s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestParse_Valid(t *testing.T) {
	tests := []struct {
		from, to string
	}{
		{"2024-01-01", "2024-12-31"},
		{"2024-01-01T10:00:00", ""},
		{"", "2024-06-01T00:00:00Z"},
		{"", ""},
	}
	for _, tt := range tests {
		_, err := Parse(tt.from, tt.to)
		if err != nil {
			t.Errorf("Parse(%q, %q) unexpected error: %v", tt.from, tt.to, err)
		}
	}
}

func TestParse_Invalid(t *testing.T) {
	tests := []struct{ from, to string }{
		{"not-a-date", ""},
		{"", "bad"},
		{"2024-12-31", "2024-01-01"}, // to before from
	}
	for _, tt := range tests {
		_, err := Parse(tt.from, tt.to)
		if err == nil {
			t.Errorf("Parse(%q, %q) expected error, got nil", tt.from, tt.to)
		}
	}
}

func TestMatch_WithinRange(t *testing.T) {
	from := mustTime("2024-03-01")
	to := mustTime("2024-03-31")
	r := &Range{From: &from, To: &to}

	entry := map[string]any{"ts": "2024-03-15T12:00:00"}
	if !r.Match(entry, "ts") {
		t.Error("expected match for timestamp within range")
	}
}

func TestMatch_OutsideRange(t *testing.T) {
	from := mustTime("2024-03-01")
	to := mustTime("2024-03-31")
	r := &Range{From: &from, To: &to}

	entry := map[string]any{"ts": "2024-04-01T00:00:00"}
	if r.Match(entry, "ts") {
		t.Error("expected no match for timestamp outside range")
	}
}

func TestMatch_MissingField(t *testing.T) {
	from := mustTime("2024-01-01")
	r := &Range{From: &from}
	entry := map[string]any{"level": "info"}
	if r.Match(entry, "ts") {
		t.Error("expected no match when timestamp field is absent")
	}
}

func TestMatch_NoRange(t *testing.T) {
	r := &Range{}
	entry := map[string]any{"ts": "2024-01-01"}
	if !r.Match(entry, "ts") {
		t.Error("expected match when no bounds are set")
	}
}
