// Package timerange provides time-window filtering for structured log entries.
//
// It parses --from and --to command-line values in common timestamp formats
// (RFC3339, date-only, or datetime without timezone) and exposes a Match
// method that checks whether a log entry's timestamp field falls within the
// specified window.
package timerange
