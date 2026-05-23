package timerange

import (
	"fmt"
	"time"
)

// Range represents an optional time window with a start and/or end bound.
type Range struct {
	From *time.Time
	To   *time.Time
}

// Parse parses --from and --to flag values into a Range.
// Accepted formats: RFC3339, 2006-01-02, 2006-01-02T15:04:05.
var formats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02",
}

func Parse(from, to string) (*Range, error) {
	r := &Range{}
	if from != "" {
		t, err := parseTime(from)
		if err != nil {
			return nil, fmt.Errorf("invalid --from value %q: %w", from, err)
		}
		r.From = &t
	}
	if to != "" {
		t, err := parseTime(to)
		if err != nil {
			return nil, fmt.Errorf("invalid --to value %q: %w", to, err)
		}
		r.To = &t
	}
	if r.From != nil && r.To != nil && r.To.Before(*r.From) {
		return nil, fmt.Errorf("--to must not be before --from")
	}
	return r, nil
}

func parseTime(s string) (time.Time, error) {
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised time format")
}

// Match reports whether the given log entry timestamp falls within the range.
// tsField is the JSON field name that holds the timestamp string.
func (r *Range) Match(entry map[string]any, tsField string) bool {
	if r.From == nil && r.To == nil {
		return true
	}
	v, ok := entry[tsField]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	t, err := parseTime(s)
	if err != nil {
		return false
	}
	if r.From != nil && t.Before(*r.From) {
		return false
	}
	if r.To != nil && t.After(*r.To) {
		return false
	}
	return true
}
