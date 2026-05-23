// Package limit provides a simple record limiter for capping output.
package limit

import "errors"

// ErrLimitReached is returned by Add when the limit has been reached.
var ErrLimitReached = errors.New("limit reached")

// Limiter tracks how many records have been processed and enforces a maximum.
// A limit of zero means unlimited.
type Limiter struct {
	max   int
	count int
}

// New creates a Limiter with the given maximum. A max of 0 means unlimited.
func New(max int) (*Limiter, error) {
	if max < 0 {
		return nil, errors.New("limit must be non-negative (0 = unlimited)")
	}
	return &Limiter{max: max}, nil
}

// Add records one more item. It returns ErrLimitReached if the limit has been
// hit, or nil if the item should be passed through.
func (l *Limiter) Add() error {
	if l.max == 0 {
		return nil
	}
	if l.count >= l.max {
		return ErrLimitReached
	}
	l.count++
	return nil
}

// Count returns the number of items accepted so far.
func (l *Limiter) Count() int {
	return l.count
}

// Done reports whether the limit has been reached.
func (l *Limiter) Done() bool {
	if l.max == 0 {
		return false
	}
	return l.count >= l.max
}

// Reset resets the counter to zero.
func (l *Limiter) Reset() {
	l.count = 0
}
