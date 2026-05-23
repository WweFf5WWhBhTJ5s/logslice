package highlight

import (
	"strings"
	"testing"
)

func TestNew_Enabled(t *testing.T) {
	h := New(true)
	if h == nil {
		t.Fatal("expected non-nil Highlighter")
	}
	if !h.enabled {
		t.Error("expected enabled=true")
	}
}

func TestNew_Disabled(t *testing.T) {
	h := New(false)
	if h.enabled {
		t.Error("expected enabled=false")
	}
}

func TestLevel_KnownLevels(t *testing.T) {
	h := New(true)
	cases := []struct {
		level string
		color string
	}{
		{"info", Green},
		{"debug", Cyan},
		{"warn", Yellow},
		{"error", Red},
		{"fatal", Red + Bold},
	}
	for _, tc := range cases {
		out := h.Level(tc.level)
		if !strings.Contains(out, tc.color) {
			t.Errorf("Level(%q): expected color %q in output %q", tc.level, tc.color, out)
		}
		if !strings.Contains(out, tc.level) {
			t.Errorf("Level(%q): expected level text in output %q", tc.level, out)
		}
		if !strings.HasSuffix(out, Reset) {
			t.Errorf("Level(%q): expected Reset suffix in output %q", tc.level, out)
		}
	}
}

func TestLevel_UnknownLevel(t *testing.T) {
	h := New(true)
	out := h.Level("trace")
	if !strings.Contains(out, "trace") {
		t.Errorf("expected 'trace' in output, got %q", out)
	}
}

func TestLevel_Disabled(t *testing.T) {
	h := New(false)
	out := h.Level("error")
	if out != "error" {
		t.Errorf("disabled highlighter should return plain text, got %q", out)
	}
}

func TestKey_Enabled(t *testing.T) {
	h := New(true)
	out := h.Key("timestamp")
	if !strings.Contains(out, Blue) {
		t.Errorf("expected Blue color in key output, got %q", out)
	}
	if !strings.Contains(out, "timestamp") {
		t.Errorf("expected key text in output, got %q", out)
	}
}

func TestKey_Disabled(t *testing.T) {
	h := New(false)
	out := h.Key("timestamp")
	if out != "timestamp" {
		t.Errorf("expected plain text, got %q", out)
	}
}

func TestValue_Enabled(t *testing.T) {
	h := New(true)
	out := h.Value("hello world")
	if !strings.Contains(out, Cyan) {
		t.Errorf("expected Cyan in value output, got %q", out)
	}
}

func TestHighlight_Custom(t *testing.T) {
	h := New(true)
	out := h.Highlight(Bold, "important")
	if !strings.HasPrefix(out, Bold) {
		t.Errorf("expected Bold prefix, got %q", out)
	}
	if !strings.HasSuffix(out, Reset) {
		t.Errorf("expected Reset suffix, got %q", out)
	}
}
