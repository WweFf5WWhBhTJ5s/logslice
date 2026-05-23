package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// Op represents a comparison operator.
type Op string

const (
	OpEq  Op = "="
	OpNeq Op = "!="
	OpGt  Op = ">"
	OpLt  Op = "<"
	OpGte Op = ">="
	OpLte Op = "<="
)

// Filter represents a single field comparison filter.
type Filter struct {
	Field string
	Op    Op
	Value string
}

// Match reports whether the given log entry (as a map) satisfies the filter.
func (f *Filter) Match(entry map[string]interface{}) bool {
	val, ok := entry[f.Field]
	if !ok {
		return false
	}

	switch v := val.(type) {
	case string:
		return compareString(v, f.Op, f.Value)
	case float64:
		right, err := strconv.ParseFloat(f.Value, 64)
		if err != nil {
			return false
		}
		return compareFloat(v, f.Op, right)
	case bool:
		return fmt.Sprintf("%v", v) == f.Value && f.Op == OpEq
	}
	return false
}

func compareString(left string, op Op, right string) bool {
	switch op {
	case OpEq:
		return left == right
	case OpNeq:
		return left != right
	case OpGt:
		return left > right
	case OpLt:
		return left < right
	case OpGte:
		return left >= right
	case OpLte:
		return left <= right
	}
	return false
}

func compareFloat(left float64, op Op, right float64) bool {
	switch op {
	case OpEq:
		return left == right
	case OpNeq:
		return left != right
	case OpGt:
		return left > right
	case OpLt:
		return left < right
	case OpGte:
		return left >= right
	case OpLte:
		return left <= right
	}
	return false
}

// Parse parses a filter expression like "level=error" or "status>=400".
func Parse(expr string) (*Filter, error) {
	ops := []Op{OpGte, OpLte, OpNeq, OpGt, OpLt, OpEq}
	for _, op := range ops {
		if idx := strings.Index(expr, string(op)); idx > 0 {
			field := strings.TrimSpace(expr[:idx])
			value := strings.TrimSpace(expr[idx+len(op):])
			if field == "" || value == "" {
				return nil, fmt.Errorf("invalid filter expression: %q", expr)
			}
			return &Filter{Field: field, Op: op, Value: value}, nil
		}
	}
	return nil, fmt.Errorf("no valid operator found in expression: %q", expr)
}
