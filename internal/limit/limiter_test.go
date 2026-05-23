package limit_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/limit"
)

func TestNew_Valid(t *testing.T) {
	for _, max := range []int{0, 1, 100} {
		l, err := limit.New(max)
		if err != nil {
			t.Fatalf("New(%d) unexpected error: %v", max, err)
		}
		if l == nil {
			t.Fatalf("New(%d) returned nil", max)
		}
	}
}

func TestNew_Negative(t *testing.T) {
	_, err := limit.New(-1)
	if err == nil {
		t.Fatal("expected error for negative limit, got nil")
	}
}

func TestLimiter_Unlimited(t *testing.T) {
	l, _ := limit.New(0)
	for i := 0; i < 1000; i++ {
		if err := l.Add(); err != nil {
			t.Fatalf("Add() returned error for unlimited limiter at i=%d: %v", i, err)
		}
	}
	if l.Done() {
		t.Fatal("Done() should be false for unlimited limiter")
	}
}

func TestLimiter_EnforcesMax(t *testing.T) {
	l, _ := limit.New(3)

	for i := 0; i < 3; i++ {
		if err := l.Add(); err != nil {
			t.Fatalf("Add() unexpected error on call %d: %v", i+1, err)
		}
	}

	if err := l.Add(); err != limit.ErrLimitReached {
		t.Fatalf("expected ErrLimitReached, got %v", err)
	}
}

func TestLimiter_Done(t *testing.T) {
	l, _ := limit.New(2)
	if l.Done() {
		t.Fatal("Done() should be false before any adds")
	}
	l.Add() //nolint:errcheck
	if l.Done() {
		t.Fatal("Done() should be false after one add with max=2")
	}
	l.Add() //nolint:errcheck
	if !l.Done() {
		t.Fatal("Done() should be true after reaching max")
	}
}

func TestLimiter_Count(t *testing.T) {
	l, _ := limit.New(5)
	for i := 1; i <= 3; i++ {
		l.Add() //nolint:errcheck
		if l.Count() != i {
			t.Fatalf("Count() = %d, want %d", l.Count(), i)
		}
	}
}

func TestLimiter_Reset(t *testing.T) {
	l, _ := limit.New(2)
	l.Add() //nolint:errcheck
	l.Add() //nolint:errcheck
	l.Reset()
	if l.Count() != 0 {
		t.Fatalf("Count() after Reset = %d, want 0", l.Count())
	}
	if l.Done() {
		t.Fatal("Done() should be false after Reset")
	}
}
