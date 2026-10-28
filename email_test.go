package main

import (
	"strings"
	"testing"
	"time"
)

func TestSubjectFor(t *testing.T) {
	day := time.Date(2026, 10, 21, 5, 17, 0, 0, time.UTC)
	if got, want := subjectFor(day), "Morning texts: Wed, Oct 21"; got != want {
		t.Errorf("subjectFor = %q, want %q", got, want)
	}
}

func TestFormatBody(t *testing.T) {
	body := formatBody(buildDrafts(testQuotes(), goodLines()))
	for _, want := range []string{
		"OPTION 1\n\n\"Pressure is a privilege.\"\n- Billie Jean King\n\nThe hard call today",
		"OPTION 3\n\n\"The best way out is always through.\"\n- Robert Frost\n",
		"OPTION 4\n\n",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q:\n%s", want, body)
		}
	}
	if got := strings.Count(body, "\n-----\n"); got != 3 {
		t.Errorf("got %d dividers, want 3", got)
	}
}
