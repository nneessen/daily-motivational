package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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

func TestBuildPromptListsQuotesInOrder(t *testing.T) {
	day := time.Date(2026, 10, 29, 5, 17, 0, 0, time.UTC)
	got := buildPrompt(day, testQuotes(), "Contest ends Friday.")
	for _, want := range []string{
		"Today is Thursday, October 29.",
		"Contest ends Friday.",
		"This morning's 4 quotes",
		"1. \"Pressure is a privilege.\" - Billie Jean King",
		"4. \"Well done is better than well said.\" - Benjamin Franklin",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("prompt missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(buildPrompt(day, testQuotes(), ""), "going on with the team") {
		t.Error("empty focus should leave the team notes out")
	}
}
