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
// each morning so the options don't all say the same thing. Most are about mindset,
// not the job; the last few tie back to the work.
var angles = []string{
	"whatever it takes: once you're really committed, you stop asking if and start asking how",
	"discipline beats mood: you don't wait to feel like it, you start and the feeling catches up",
	"confidence is built from evidence: reps and proof, not hyping yourself up or waiting to feel ready",
	"identity first: decide who you are, and the actions follow before the results show up",
	"a no is information, not a verdict on you",
	"what you say to yourself right after a loss decides the next hour",
	"control what you control: your effort, your attitude, your activity, and let the rest go",
	"nerves and excitement feel the same in your body, so call it excitement and use it",
	"the thing you keep putting off is usually the thing that moves you forward, so do it first",
	"skill gets built in boring reps nobody sees, and it stacks up faster than people expect",
	"ownership: the second you blame the leads, the market or the day, you hand your power away",
	"your real standard is what you do on your worst day, not your best",
	"stay hungry after a win: a good week is exactly when most people ease off",
	"the people you spend time with and what you put in your head set your ceiling",
	"a specific number you look at every day beats a vague wish",
	"the work you put in now pays you months from now, so keep stacking even when you can't see it yet",
	"the fear of looking dumb holds people back more than any lack of talent",
	"be coachable: the people who grow fastest take feedback without getting defensive",
	"your energy is a decision you make before the day starts, and the room feels it",
	"keep the promises you make to yourself: each one builds trust in yourself, each broken one chips it away",
	"decide ahead of time what you'll do when the hard moment shows up, so you don't negotiate with yourself in it",
	"repetition and boredom are where most people quit, and that's where the separation happens",
	"playing it safe feels free, but it costs you slowly in everything you didn't go after",
	"who's counting on you at home is the reason to push when it gets heavy",
	"every person who calls already reached out, and one more question might protect their whole family",
	"one conversation can cover a whole household: a spouse, grown kids, parents who never got set up",
	"income is a scoreboard for how many families you've protected, so chase the families",
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
	b.WriteString("Background only. Don't build the messages around the day or the time of year. At most one option may mention either.\n")
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
