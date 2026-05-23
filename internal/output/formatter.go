package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format defines the output format type.
type Format string

const (
	FormatJSON    Format = "json"
	FormatPretty  Format = "pretty"
	FormatCompact Format = "compact"
)

// Formatter writes log entries to an output writer.
type Formatter struct {
	writer io.Writer
	format Format
	fields []string
}

// New creates a new Formatter with the given writer, format, and optional field selection.
func New(w io.Writer, format Format, fields []string) *Formatter {
	return &Formatter{
		writer: w,
		format: format,
		fields: fields,
	}
}

// Write outputs a single log entry according to the configured format.
func (f *Formatter) Write(entry map[string]interface{}) error {
	data := f.selectFields(entry)

	switch f.format {
	case FormatPretty:
		return f.writePretty(data)
	case FormatCompact:
		return f.writeCompact(data)
	default:
		return f.writeJSON(data)
	}
}

func (f *Formatter) selectFields(entry map[string]interface{}) map[string]interface{} {
	if len(f.fields) == 0 {
		return entry
	}
	selected := make(map[string]interface{}, len(f.fields))
	for _, field := range f.fields {
		if val, ok := entry[field]; ok {
			selected[field] = val
		}
	}
	return selected
}

func (f *Formatter) writeJSON(data map[string]interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("output: marshal error: %w", err)
	}
	_, err = fmt.Fprintln(f.writer, string(b))
	return err
}

func (f *Formatter) writePretty(data map[string]interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("output: marshal error: %w", err)
	}
	_, err = fmt.Fprintln(f.writer, string(b))
	return err
}

func (f *Formatter) writeCompact(data map[string]interface{}) error {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(data))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, data[k]))
	}
	_, err := fmt.Fprintln(f.writer, strings.Join(parts, " "))
	return err
}
