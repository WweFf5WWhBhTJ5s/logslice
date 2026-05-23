// Package count provides a simple entry counter for logslice pipelines.
// It tracks the total number of log entries processed and can report
// the final count to any io.Writer.
package count

import (
	"fmt"
	"io"
	"sync/atomic"
)

// Counter tracks the number of log entries that have passed through the
// pipeline. It is safe for concurrent use.
type Counter struct {
	total atomic.Int64
	enabled bool
}

// New returns a new Counter. When enabled is false, Increment is a no-op
// and Report writes nothing, allowing callers to disable counting with zero
// overhead.
func New(enabled bool) *Counter {
	return &Counter{enabled: enabled}
}

// Increment records one additional entry. It is safe to call from multiple
// goroutines.
func (c *Counter) Increment() {
	if !c.enabled {
		return
	}
	c.total.Add(1)
}

// Total returns the current count.
func (c *Counter) Total() int64 {
	return c.total.Load()
}

// Report writes a human-readable summary line to w.
// If the counter is disabled it writes nothing and returns nil.
func (c *Counter) Report(w io.Writer) error {
	if !c.enabled {
		return nil
	}
	_, err := fmt.Fprintf(w, "matched %d entries\n", c.total.Load())
	return err
}
