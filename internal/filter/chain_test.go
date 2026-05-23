package filter

import "testing"

func TestNewChain_Valid(t *testing.T) {
	c, err := NewChain([]string{"level=error", "status>=500"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 2 {
		t.Errorf("expected 2 filters, got %d", c.Len())
	}
}

func TestNewChain_Invalid(t *testing.T) {
	_, err := NewChain([]string{"level=error", "badexpr"})
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestChain_Match_AllPass(t *testing.T) {
	c, _ := NewChain([]string{"level=error", "status>=500"})
	entry := map[string]interface{}{
		"level":  "error",
		"status": float64(503),
	}
	if !c.Match(entry) {
		t.Error("expected chain to match")
	}
}

func TestChain_Match_OneFails(t *testing.T) {
	c, _ := NewChain([]string{"level=error", "status>=500"})
	entry := map[string]interface{}{
		"level":  "error",
		"status": float64(200),
	}
	if c.Match(entry) {
		t.Error("expected chain not to match")
	}
}

func TestChain_Match_Empty(t *testing.T) {
	c, _ := NewChain([]string{})
	entry := map[string]interface{}{"anything": "value"}
	if !c.Match(entry) {
		t.Error("empty chain should always match")
	}
}
