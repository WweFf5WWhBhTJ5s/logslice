package dedupe

import (
	"testing"
)

func TestNew_Valid(t *testing.T) {
	d, err := New("request_id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Deduplicator")
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestIsDuplicate_FirstOccurrence(t *testing.T) {
	d, _ := New("id")
	entry := map[string]interface{}{"id": "abc123", "msg": "hello"}
	if d.IsDuplicate(entry) {
		t.Error("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicate_SecondOccurrence(t *testing.T) {
	d, _ := New("id")
	entry := map[string]interface{}{"id": "abc123"}
	d.IsDuplicate(entry) // first — mark as seen
	if !d.IsDuplicate(entry) {
		t.Error("second occurrence should be a duplicate")
	}
}

func TestIsDuplicate_MissingField(t *testing.T) {
	d, _ := New("id")
	entry := map[string]interface{}{"msg": "no id here"}
	if d.IsDuplicate(entry) {
		t.Error("entry missing field should never be a duplicate")
	}
	if d.IsDuplicate(entry) {
		t.Error("repeated entry missing field should still not be a duplicate")
	}
}

func TestIsDuplicate_DifferentValues(t *testing.T) {
	d, _ := New("id")
	e1 := map[string]interface{}{"id": "aaa"}
	e2 := map[string]interface{}{"id": "bbb"}
	if d.IsDuplicate(e1) || d.IsDuplicate(e2) {
		t.Error("distinct values should not be duplicates")
	}
	if d.SeenCount() != 2 {
		t.Errorf("expected SeenCount 2, got %d", d.SeenCount())
	}
}

func TestReset(t *testing.T) {
	d, _ := New("id")
	entry := map[string]interface{}{"id": "x"}
	d.IsDuplicate(entry)
	d.Reset()
	if d.IsDuplicate(entry) {
		t.Error("after reset, entry should not be considered a duplicate")
	}
	if d.SeenCount() != 1 {
		t.Errorf("expected SeenCount 1 after reset+one entry, got %d", d.SeenCount())
	}
}
