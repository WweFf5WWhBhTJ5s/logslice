package count_test

import (
	"bytes"
	"testing"

	"github.com/logslice/logslice/internal/count"
)

func TestNew_Enabled(t *testing.T) {
	c := count.New(true)
	if c == nil {
		t.Fatal("expected non-nil Counter")
	}
}

func TestNew_Disabled(t *testing.T) {
	c := count.New(false)
	if c == nil {
		t.Fatal("expected non-nil Counter")
	}
}

func TestCounter_IncrementAndTotal(t *testing.T) {
	c := count.New(true)
	for i := 0; i < 5; i++ {
		c.Increment()
	}
	if got := c.Total(); got != 5 {
		t.Fatalf("expected total 5, got %d", got)
	}
}

func TestCounter_Disabled_NoOp(t *testing.T) {
	c := count.New(false)
	for i := 0; i < 10; i++ {
		c.Increment()
	}
	if got := c.Total(); got != 0 {
		t.Fatalf("expected total 0 when disabled, got %d", got)
	}
}

func TestCounter_Report_Enabled(t *testing.T) {
	c := count.New(true)
	c.Increment()
	c.Increment()
	c.Increment()

	var buf bytes.Buffer
	if err := c.Report(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "matched 3 entries\n"
	if got := buf.String(); got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestCounter_Report_Disabled(t *testing.T) {
	c := count.New(false)
	c.Increment()

	var buf bytes.Buffer
	if err := c.Report(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "" {
		t.Fatalf("expected empty output when disabled, got %q", got)
	}
}

func TestCounter_ZeroByDefault(t *testing.T) {
	c := count.New(true)
	if got := c.Total(); got != 0 {
		t.Fatalf("expected initial total 0, got %d", got)
	}
}
