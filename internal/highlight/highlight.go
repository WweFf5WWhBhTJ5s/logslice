package highlight

import (
	"fmt"
	"strings"
)

// Color ANSI codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// LevelColors maps log level strings to ANSI color codes.
var LevelColors = map[string]string{
	"debug": Cyan,
	"info":  Green,
	"warn":  Yellow,
	"error": Red,
	"fatal": Red + Bold,
}

// Highlighter applies ANSI color codes to log output.
type Highlighter struct {
	enabled bool
}

// New returns a Highlighter. If enabled is false, all methods return
// plain text without any ANSI escape sequences.
func New(enabled bool) *Highlighter {
	return &Highlighter{enabled: enabled}
}

// Level colorizes a log level string based on its severity.
func (h *Highlighter) Level(level string) string {
	if !h.enabled {
		return level
	}
	color, ok := LevelColors[strings.ToLower(level)]
	if !ok {
		color = Reset
	}
	return fmt.Sprintf("%s%s%s", color, level, Reset)
}

// Key colorizes a JSON field key.
func (h *Highlighter) Key(key string) string {
	if !h.enabled {
		return key
	}
	return fmt.Sprintf("%s%s%s", Blue, key, Reset)
}

// Value colorizes a JSON field value string.
func (h *Highlighter) Value(value string) string {
	if !h.enabled {
		return value
	}
	return fmt.Sprintf("%s%s%s", Cyan, value, Reset)
}

// Highlight wraps arbitrary text with a given ANSI color code.
func (h *Highlighter) Highlight(color, text string) string {
	if !h.enabled {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}
