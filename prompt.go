package main

import (
	_ "embed"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

// voiceGuide is the system prompt. Edit voice.md to change the tone.
//
//go:embed voice.md
var voiceGuide string

// angles are the ideas a message can be built around. A different four are picked
// each morning so the options don't all say the same thing.
var angles = []string{
	"take every call, no cherry-picking: the billing and 'quick question' calls count",
	"listen more than you talk and let the caller tell you what they need",
	"ask one more question than feels natural",
	"the whole household: spouse, grown kids starting families, parents who never got set up",
	"bouncing back after a rough call or a slow day",
	"consistency over hype: same effort on the last call as the first",
	"the mission: every policy written is a family protected",
	"trust the system, it works when you work it",
	"the first call of the day sets the tone",
	"finish the day strong, the last hour counts the same",
	"make an attempt on every call, even when it feels like a long shot",
	"stay present: callers can hear when you're distracted",
	"get curious about the caller's situation",
	"no two calls are the same, so reset after each one",
}

func pickAngles(n int) []string {
	picked := make([]string, 0, n)
	for _, i := range rand.Perm(len(angles))[:n] {
		picked = append(picked, angles[i])
	}
	return picked
}

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

func weekdayNote(day time.Weekday) string {
	switch day {
	case time.Monday:
		return "It's the start of the week."
	case time.Friday:
		return "It's the last day of the work week."
	default:
		return "It's midweek."
	}
}

// calendarNotes says where the date falls in the month, quarter and year.
func calendarNotes(t time.Time) []string {
	var notes []string
	lastDay := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
	daysLeft := lastDay - t.Day()
	switch {
	case daysLeft == 0:
		notes = append(notes, "It's the last day of the month.")
	case daysLeft <= 3:
		notes = append(notes, "It's the last few days of the month.")
	case t.Day() <= 3:
		notes = append(notes, "It's the first few days of a new month.")
	}
	if daysLeft <= 3 && t.Month()%3 == 0 {
		notes = append(notes, "The quarter ends with this month.")
	}
	if t.Month() >= time.October {
		notes = append(notes, "It's Q4, the final stretch of the year.")
	}
	return notes
}

func buildPrompt(today time.Time, picked []string, focus string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Today is %s. %s\n", today.Format("Monday, January 2, 2006"), weekdayNote(today.Weekday()))
	for _, note := range calendarNotes(today) {
		b.WriteString(note + "\n")
	}
	if focus != "" {
		fmt.Fprintf(&b, "\nWhat's going on with the team right now (work it in where it fits, don't force it):\n%s\n", focus)
	}
	b.WriteString("\nWrite this morning's options, each built around a different angle:\n")
	for i, angle := range picked {
		size := "short"
		if i >= shortCount {
			size = "longer"
		}
		fmt.Fprintf(&b, "%d. (%s) %s\n", i+1, size, angle)
	}
	return b.String()
}
