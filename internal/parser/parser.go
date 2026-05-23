package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// LogEntry represents a single parsed JSON log line.
type LogEntry struct {
	Raw    string
	Fields map[string]interface{}
}

// Parser reads log lines from a reader and parses them as JSON.
type Parser struct {
	reader  io.Reader
	scanner *bufio.Scanner
}

// New creates a new Parser that reads from r.
func New(r io.Reader) *Parser {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &Parser{
		reader:  r,
		scanner: scanner,
	}
}

// Next advances to the next log entry. Returns (entry, true) on success,
// (nil, false) when there are no more lines, and an error on parse failure.
func (p *Parser) Next() (*LogEntry, bool, error) {
	for p.scanner.Scan() {
		line := p.scanner.Text()
		if len(line) == 0 {
			continue
		}
		var fields map[string]interface{}
		if err := json.Unmarshal([]byte(line), &fields); err != nil {
			return nil, false, fmt.Errorf("parse error: %w (line: %q)", err, line)
		}
		return &LogEntry{Raw: line, Fields: fields}, true, nil
	}
	if err := p.scanner.Err(); err != nil {
		return nil, false, fmt.Errorf("scanner error: %w", err)
	}
	return nil, false, nil
}

// All reads all log entries from the parser, skipping non-JSON lines when
// skipInvalid is true, otherwise returning on first error.
func (p *Parser) All(skipInvalid bool) ([]*LogEntry, error) {
	var entries []*LogEntry
	for {
		entry, ok, err := p.Next()
		if err != nil {
			if skipInvalid {
				continue
			}
			return nil, err
		}
		if !ok {
			break
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
