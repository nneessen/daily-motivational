package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	draftCount = 4
	shortCount = 3 // the first three options are short texts, the last one is longer

	// A little looser than the limits in voice.md, so a few extra words don't cost a retry.
	maxShortWords = 80
	minLongWords  = 90
	maxLongWords  = 280
)

// bannedPhrases make a text read like it came from a bot or a break room poster.
// Keep this in sync with the list in voice.md.
var bannedPhrases = []string{
	"game-changer", "game changer", "unlock", "elevate", "journey", "crush it",
	"embrace", "harness", "level up", "synergy", "rockstar", "superstar",
	"ignite", "unleash", "dive in", "let that sink in", "read that again", "you got this",
	"you've got this", "at the end of the day", "here's the thing", "in today's",
	"good morning team", "rise and shine", "let's talk about", "i get it", "we've all been there",
	// gym-wall lines
	"comfort zone", "1% better", "rent is due", "the grind", "hard work beats talent",
	"not a sprint", "no pain no gain", "hardest worker in the room", "only competition is yourself",
}

// calendarTalk spots a message built around the day or the time of year.
var calendarTalk = regexp.MustCompile(`(?i)\b(monday|tuesday|wednesday|thursday|friday|weekend|q4|quarter|final stretch|end of the year)\b`)

// wifeOrHusband guesses who's calling. Nick says "spouse".
var wifeOrHusband = regexp.MustCompile(`(?i)\b(wife|wives|husbands?)\b`)

var (
	hashtag      = regexp.MustCompile(`(^|\s)#[A-Za-z]`)
	emDash       = regexp.MustCompile(`[ \t]*—[ \t]*`)
	spacedEnDash = regexp.MustCompile(`[ \t]+–[ \t]*`)
)

// findProblems lists every way the drafts break the voice rules. Empty means good to send.
func findProblems(drafts []Draft) []string {
	var problems []string
	if len(drafts) != draftCount {
		problems = append(problems, fmt.Sprintf("there were %d options instead of %d", len(drafts), draftCount))
	}
	calendarCount := 0
	for _, d := range drafts {
		if calendarTalk.MatchString(d.Text) {
			calendarCount++
		}
	}
	if calendarCount > 1 {
		problems = append(problems, fmt.Sprintf("%d options mention the day or the time of year, only one may", calendarCount))
	}
	for i, d := range drafts {
		n := i + 1
		lower := strings.ToLower(d.Text)
		if strings.ContainsAny(d.Text, "—–") {
			problems = append(problems, fmt.Sprintf("option %d uses a dash (— or –)", n))
		}
		for _, phrase := range bannedPhrases {
			if strings.Contains(lower, phrase) {
				problems = append(problems, fmt.Sprintf("option %d says %q", n, phrase))
			}
		}
		if wifeOrHusband.MatchString(d.Text) {
			problems = append(problems, fmt.Sprintf("option %d says wife or husband instead of spouse", n))
		}
		if hashtag.MatchString(d.Text) {
			problems = append(problems, fmt.Sprintf("option %d has a hashtag", n))
		}
		if emojiCount(d.Text) > 1 {
			problems = append(problems, fmt.Sprintf("option %d has more than one emoji", n))
		}
		if strings.Count(d.Text, "!") > 2 {
			problems = append(problems, fmt.Sprintf("option %d has more than two exclamation points", n))
		}
		if strings.Contains(d.Text, "**") {
			problems = append(problems, fmt.Sprintf("option %d uses markdown", n))
		}

		words := len(strings.Fields(d.Text))
		if i < shortCount && words > maxShortWords {
			problems = append(problems, fmt.Sprintf("option %d is %d words, short ones stay under 70", n, words))
		}
		if i >= shortCount && (words < minLongWords || words > maxLongWords) {
			problems = append(problems, fmt.Sprintf("option %d is %d words, the longer one should be 100 to 250", n, words))
		}
	}
	return problems
}

func emojiCount(s string) int {
	count := 0
	for _, r := range s {
		if unicode.Is(unicode.So, r) {
			count++
		}
	}
	return count
}

func retryNote(problems []string) string {
	return "\nA first draft for today came back with these problems. Write all four fresh and avoid them:\n- " +
		strings.Join(problems, "\n- ") + "\n"
}

// cleanText swaps dashes for the "..." pause Nick actually types. En dashes inside
// ranges like 2–5 become plain hyphens.
func cleanText(s string) string {
	s = emDash.ReplaceAllString(s, "... ")
	s = spacedEnDash.ReplaceAllString(s, "... ")
	s = strings.ReplaceAll(s, "–", "-")
	return strings.ReplaceAll(s, "... \n", "...\n")
}

func cleanDrafts(drafts []Draft) []Draft {
	cleaned := make([]Draft, len(drafts))
	for i, d := range drafts {
		cleaned[i] = Draft{Label: d.Label, Text: strings.TrimSpace(cleanText(d.Text))}
	}
	return cleaned
}
