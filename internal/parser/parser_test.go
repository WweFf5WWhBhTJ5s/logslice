package parser

import (
	"strings"
	"testing"
)

func TestNext_ValidJSON(t *testing.T) {
	input := `{"level":"info","msg":"started"}
{"level":"error","msg":"failed","code":500}
`
	p := New(strings.NewReader(input))

	entry, ok, err := p.Next()
	if err != nil || !ok {
		t.Fatalf("expected first entry, got err=%v ok=%v", err, ok)
	}
	if entry.Fields["level"] != "info" {
		t.Errorf("expected level=info, got %v", entry.Fields["level"])
	}

	entry, ok, err = p.Next()
	if err != nil || !ok {
		t.Fatalf("expected second entry, got err=%v ok=%v", err, ok)
	}
	if entry.Fields["msg"] != "failed" {
		t.Errorf("expected msg=failed, got %v", entry.Fields["msg"])
	}

	_, ok, err = p.Next()
	if err != nil || ok {
		t.Errorf("expected EOF, got err=%v ok=%v", err, ok)
	}
}

func TestNext_InvalidJSON(t *testing.T) {
	input := `not json`
	p := New(strings.NewReader(input))
	_, _, err := p.Next()
	if err == nil {
		t.Error("expected parse error for invalid JSON")
	}
}

func TestNext_EmptyLines(t *testing.T) {
	input := "\n\n{\"level\":\"debug\"}\n\n"
	p := New(strings.NewReader(input))
	entry, ok, err := p.Next()
	if err != nil || !ok {
		t.Fatalf("expected entry after empty lines, got err=%v ok=%v", err, ok)
	}
	if entry.Fields["level"] != "debug" {
		t.Errorf("expected level=debug, got %v", entry.Fields["level"])
	}
}

func TestAll_SkipInvalid(t *testing.T) {
	input := `{"level":"info"}
not json
{"level":"warn"}
`
	p := New(strings.NewReader(input))
	entries, err := p.All(true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestAll_StopOnInvalid(t *testing.T) {
	input := `{"level":"info"}
not json
`
	p := New(strings.NewReader(input))
	_, err := p.All(false)
	if err == nil {
		t.Error("expected error when skipInvalid is false")
	}
}
