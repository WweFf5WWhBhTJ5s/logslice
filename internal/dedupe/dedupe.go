// Package dedupe provides deduplication of log entries based on a specified field.
package dedupe

import (
	"errors"
	"fmt"
)

// Deduplicator tracks seen values for a given field and filters duplicate entries.
type Deduplicator struct {
	field string
	seen  map[string]struct{}
}

// New creates a new Deduplicator that deduplicates on the given field.
// Returns an error if field is empty.
func New(field string) (*Deduplicator, error) {
	if field == "" {
		return nil, errors.New("dedupe: field must not be empty")
	}
	return &Deduplicator{
		field: field,
		seen:  make(map[string]struct{}),
	}, nil
}

// IsDuplicate reports whether the given log entry (as a map) has already been
// seen for the configured field. If the field is missing from the entry, the
// entry is never considered a duplicate and is always passed through.
func (d *Deduplicator) IsDuplicate(entry map[string]interface{}) bool {
	val, ok := entry[d.field]
	if !ok {
		return false
	}
	key := fmt.Sprintf("%v", val)
	if _, exists := d.seen[key]; exists {
		return true
	}
	d.seen[key] = struct{}{}
	return false
}

// Reset clears all previously seen values, allowing the deduplicator to be
// reused for a new stream.
func (d *Deduplicator) Reset() {
	d.seen = make(map[string]struct{})
}

// SeenCount returns the number of distinct values observed so far.
func (d *Deduplicator) SeenCount() int {
	return len(d.seen)
}
