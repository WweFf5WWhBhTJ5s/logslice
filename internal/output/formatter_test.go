package output

import (
	"bytes"
	"strings"
	"testing"
)

func entry() map[string]interface{} {
	return map[string]interface{}{
		"level":   "info",
		"message": "hello world",
		"code":    float64(200),
	}
}

func TestFormatter_WriteJSON(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatJSON, nil)
	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"level":"info"`) {
		t.Errorf("expected level field in JSON output, got: %s", out)
	}
}

func TestFormatter_WritePretty(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatPretty, nil)
	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\n") {
		t.Errorf("expected indented output, got: %s", out)
	}
	if !strings.Contains(out, `"message": "hello world"`) {
		t.Errorf("expected message field in pretty output, got: %s", out)
	}
}

func TestFormatter_WriteCompact(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatCompact, nil)
	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if !strings.Contains(out, "level=info") {
		t.Errorf("expected level=info in compact output, got: %s", out)
	}
}

func TestFormatter_FieldSelection(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatJSON, []string{"level"})
	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "message") {
		t.Errorf("expected message to be excluded, got: %s", out)
	}
	if !strings.Contains(out, "level") {
		t.Errorf("expected level to be included, got: %s", out)
	}
}

func TestFormatter_FieldSelection_MissingField(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatJSON, []string{"nonexistent"})
	if err := f.Write(entry()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "{}" {
		t.Errorf("expected empty object for missing fields, got: %s", out)
	}
}
