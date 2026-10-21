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
	body := formatBody(goodDrafts())
	for _, want := range []string{"OPTION 1 (short): every call", "OPTION 3 (short): whole household", "OPTION 4 (longer): the mission"} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q:\n%s", want, body)
		}
	}
	if got := strings.Count(body, "\n-----\n"); got != 3 {
		t.Errorf("got %d dividers, want 3", got)
	}
}
