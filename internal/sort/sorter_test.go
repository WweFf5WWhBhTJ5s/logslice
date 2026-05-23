package sort

import (
	"testing"
)

func TestNew_Valid(t *testing.T) {
	for _, tc := range []struct{ dir string }{{""}, {"asc"}, {"desc"}} {
		s, err := New("ts", tc.dir)
		if err != nil {
			t.Fatalf("unexpected error for dir=%q: %v", tc.dir, err)
		}
		if s == nil {
			t.Fatal("expected non-nil sorter")
		}
	}
}

func TestNew_Invalid(t *testing.T) {
	if _, err := New("", "asc"); err == nil {
		t.Fatal("expected error for empty field")
	}
	if _, err := New("ts", "random"); err == nil {
		t.Fatal("expected error for invalid direction")
	}
}

func TestSort_StringAsc(t *testing.T) {
	s, _ := New("level", "asc")
	entries := []map[string]any{
		{"level": "warn"},
		{"level": "error"},
		{"level": "info"},
	}
	s.Sort(entries)
	expected := []string{"error", "info", "warn"}
	for i, e := range entries {
		if e["level"] != expected[i] {
			t.Errorf("index %d: got %v, want %v", i, e["level"], expected[i])
		}
	}
}

func TestSort_FloatDesc(t *testing.T) {
	s, _ := New("latency", "desc")
	entries := []map[string]any{
		{"latency": float64(10)},
		{"latency": float64(50)},
		{"latency": float64(30)},
	}
	s.Sort(entries)
	expected := []float64{50, 30, 10}
	for i, e := range entries {
		if e["latency"] != expected[i] {
			t.Errorf("index %d: got %v, want %v", i, e["latency"], expected[i])
		}
	}
}

func TestSort_MissingFieldSinksToBottom(t *testing.T) {
	s, _ := New("ts", "asc")
	entries := []map[string]any{
		{"ts": "2024-01-02"},
		{"msg": "no ts here"},
		{"ts": "2024-01-01"},
	}
	s.Sort(entries)
	if entries[0]["ts"] != "2024-01-01" {
		t.Errorf("expected earliest ts first, got %v", entries[0])
	}
	if _, ok := entries[2]["ts"]; ok {
		t.Errorf("expected entry with missing ts to be last")
	}
}
