package timerange

import (
	"testing"
	"time"
)

func BenchmarkMatch(b *testing.B) {
	from := mustTime("2024-01-01")
	to := mustTime("2024-12-31")
	r := &Range{From: &from, To: &to}
	entry := map[string]any{"ts": "2024-06-15T08:30:00"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Match(entry, "ts")
	}
}

func BenchmarkParseTime(b *testing.B) {
	samples := []string{
		"2024-06-15T08:30:00Z",
		"2024-06-15T08:30:00",
		"2024-06-15",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parseTime(samples[i%len(samples)])
	}
}

func BenchmarkParse(b *testing.B) {
	now := time.Now().Format(time.RFC3339)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse("2024-01-01", now)
	}
}
