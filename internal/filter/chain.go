package filter

import "fmt"

// Chain holds multiple filters that must ALL match (AND semantics).
type Chain struct {
	filters []*Filter
}

// NewChain creates a Chain from a slice of raw filter expressions.
func NewChain(exprs []string) (*Chain, error) {
	c := &Chain{}
	for _, expr := range exprs {
		f, err := Parse(expr)
		if err != nil {
			return nil, fmt.Errorf("filter parse error: %w", err)
		}
		c.filters = append(c.filters, f)
	}
	return c, nil
}

// Match returns true if the entry satisfies all filters in the chain.
// An empty chain always returns true.
func (c *Chain) Match(entry map[string]interface{}) bool {
	for _, f := range c.filters {
		if !f.Match(entry) {
			return false
		}
	}
	return true
}

// Len returns the number of filters in the chain.
func (c *Chain) Len() int {
	return len(c.filters)
}
