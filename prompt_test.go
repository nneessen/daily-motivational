package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCalendarNotes(t *testing.T) {
	tests := []struct {
		date string
		want []string
	}{
		{"2026-10-31", []string{"last day of the month", "Q4"}},
		{"2026-09-30", []string{"last day of the month", "quarter ends"}},
		{"2026-12-15", []string{"Q4"}},
		{"2026-06-10", nil},
	}
	for _, tt := range tests {
		day, _ := time.Parse("2006-01-02", tt.date)
		notes := calendarNotes(day)
		got := strings.Join(notes, " ")
		if len(notes) != len(tt.want) {
			t.Errorf("%s: got %d notes %q, want %d", tt.date, len(notes), got, len(tt.want))
		}
		for _, want := range tt.want {
			if !strings.Contains(got, want) {
				t.Errorf("%s: notes %q missing %q", tt.date, got, want)
			}
		}
	}
}

func TestPickAnglesAreDistinct(t *testing.T) {
	seen := map[string]bool{}
	for _, angle := range pickAngles(draftCount) {
		if seen[angle] {
			t.Fatalf("angle picked twice: %q", angle)
		}
		seen[angle] = true
	}
	if len(seen) != draftCount {
		t.Fatalf("got %d angles, want %d", len(seen), draftCount)
	}
}

func TestReadFocusSkipsComments(t *testing.T) {
	path := filepath.Join(t.TempDir(), "focus.txt")
	if err := os.WriteFile(path, []byte("# a comment\n\nContest ends Friday.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := readFocus(path); got != "Contest ends Friday." {
		t.Errorf("readFocus = %q", got)
	}
	if got := readFocus(filepath.Join(t.TempDir(), "missing.txt")); got != "" {
		t.Errorf("missing file should give empty focus, got %q", got)
	}
}

func TestBuildPromptLabelsSizes(t *testing.T) {
	day := time.Date(2026, 10, 21, 5, 17, 0, 0, time.UTC)
	got := buildPrompt(day, []string{"a", "b", "c", "d"}, "Contest ends Friday.")
	for _, want := range []string{"Wednesday, October 21, 2026", "1. (short) a", "4. (longer) d", "Contest ends Friday."} {
		if !strings.Contains(got, want) {
			t.Errorf("prompt missing %q:\n%s", want, got)
		}
	}
}
