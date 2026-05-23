package sort

import (
	"fmt"
	"sort"
)

// Direction represents sort order.
type Direction int

const (
	Asc  Direction = iota
	Desc Direction = iota
)

// Sorter sorts log entries by a given field.
type Sorter struct {
	field     string
	direction Direction
}

// New creates a new Sorter for the given field and direction string ("asc" or "desc").
func New(field, dir string) (*Sorter, error) {
	if field == "" {
		return nil, fmt.Errorf("sort field must not be empty")
	}
	var d Direction
	switch dir {
	case "", "asc":
		d = Asc
	case "desc":
		d = Desc
	default:
		return nil, fmt.Errorf("invalid sort direction %q: must be \"asc\" or \"desc\"", dir)
	}
	return &Sorter{field: field, direction: d}, nil
}

// Sort sorts a slice of parsed log entries (map[string]any) in-place.
// Entries missing the sort field are placed at the end.
func (s *Sorter) Sort(entries []map[string]any) {
	sort.SliceStable(entries, func(i, j int) bool {
		vi, oki := entries[i][s.field]
		vj, okj := entries[j][s.field]

		// Missing fields sink to the bottom regardless of direction.
		if !oki && !okj {
			return false
		}
		if !oki {
			return false
		}
		if !okj {
			return true
		}

		less := lessThan(vi, vj)
		if s.direction == Desc {
			return !less && vi != vj
		}
		return less
	})
}

// lessThan compares two values of potentially mixed types.
// Numeric values are compared as float64; everything else as string.
func lessThan(a, b any) bool {
	switch av := a.(type) {
	case float64:
		if bv, ok := b.(float64); ok {
			return av < bv
		}
	case string:
		if bv, ok := b.(string); ok {
			return av < bv
		}
	}
	return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
}
