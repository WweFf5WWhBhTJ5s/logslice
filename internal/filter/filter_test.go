package filter

import (
	"testing"
)

func TestParse_ValidExpressions(t *testing.T) {
	tests := []struct {
		expr  string
		field string
		op    Op
		value string
	}{
		{"level=error", "level", OpEq, "error"},
		{"status>=400", "status", OpGte, "400"},
		{"code!=200", "code", OpNeq, "200"},
		{"latency>100", "latency", OpGt, "100"},
		{"latency<=50", "latency", OpLte, "50"},
	}
	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			f, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f.Field != tt.field || f.Op != tt.op || f.Value != tt.value {
				t.Errorf("got {%s %s %s}, want {%s %s %s}", f.Field, f.Op, f.Value, tt.field, tt.op, tt.value)
			}
		})
	}
}

func TestParse_InvalidExpressions(t *testing.T) {
	exprs := []string{"", "levelonly", "=value", "field="}
	for _, expr := range exprs {
		_, err := Parse(expr)
		if err == nil {
			t.Errorf("expected error for expr %q, got nil", expr)
		}
	}
}

func TestFilter_Match_String(t *testing.T) {
	f := &Filter{Field: "level", Op: OpEq, Value: "error"}
	entry := map[string]interface{}{"level": "error", "msg": "oops"}
	if !f.Match(entry) {
		t.Error("expected match")
	}
	entry["level"] = "info"
	if f.Match(entry) {
		t.Error("expected no match")
	}
}

func TestFilter_Match_Float(t *testing.T) {
	f := &Filter{Field: "status", Op: OpGte, Value: "400"}
	if !f.Match(map[string]interface{}{"status": float64(500)}) {
		t.Error("expected match for 500 >= 400")
	}
	if f.Match(map[string]interface{}{"status": float64(200)}) {
		t.Error("expected no match for 200 >= 400")
	}
}

func TestFilter_Match_MissingField(t *testing.T) {
	f := &Filter{Field: "level", Op: OpEq, Value: "error"}
	if f.Match(map[string]interface{}{"msg": "hello"}) {
		t.Error("expected no match when field is missing")
	}
}
