package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"time"
)

// voiceGuide is the system prompt. Edit voice.md to change the tone.
//
//go:embed voice.md
var voiceGuide string

// readFocus returns Nick's notes from focus.txt, skipping # comment lines.
// A missing or empty file just means there's nothing special going on.
func readFocus(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var lines []string
	for _, line := range strings.Split(string(raw), "\n") {
		if line = strings.TrimSpace(line); line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}

func buildPrompt(today time.Time, quotes []Quote, focus string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Today is %s. That's background, not something to build the lines around.\n\n", today.Format("Monday, January 2"))
	if focus != "" {
		fmt.Fprintf(&b, "What's going on with the team right now (use it in one line at most, only if it fits):\n%s\n\n", focus)
	}
	fmt.Fprintf(&b, "This morning's %d quotes. Write the line Nick adds under each one, in the same order:\n", len(quotes))
	for i, q := range quotes {
		fmt.Fprintf(&b, "%d. \"%s\" - %s\n", i+1, q.Text, q.Author)
	}
	return b.String()
}
