package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const (
	draftCount = 4

	// A little looser than the limit in voice.md, so a few extra words don't cost a retry.
	maxLineWords = 35
)

// bannedPhrases make a line read like it came from a bot or a break room poster.
// Keep this in sync with the list in voice.md. They're only checked in the lines
// Claude writes, never in the quotes.
var bannedPhrases = []string{
	"game-changer", "game changer", "unlock", "elevate", "journey", "crush it",
	"embrace", "harness", "level up", "synergy", "rockstar", "superstar",
	"ignite", "unleash", "dive in", "let that sink in", "read that again", "read that twice",
	"you got this", "you've got this", "at the end of the day", "here's the thing",
	"good morning team", "rise and shine", "i get it", "we've all been there",
	"in other words", "translation:", "the lesson", "this is your sign", "reminder:",
	// gym-wall lines
	"comfort zone", "1% better", "rent is due", "the grind", "hard work beats talent",
	"not a sprint", "no pain no gain", "hardest worker in the room", "only competition is yourself",
}

// calendarTalk spots a line built around the day or the time of year.
var calendarTalk = regexp.MustCompile(`(?i)\b(monday|tuesday|wednesday|thursday|friday|weekend|q4|quarter|final stretch|end of the year)\b`)

// callTalk spots a line about the phones. Only one line a morning may tie back to the job.
var callTalk = regexp.MustCompile(`(?i)\b(calls?|callers?|phones?)\b`)

// wifeOrHusband guesses who's calling. Nick says "spouse".
var wifeOrHusband = regexp.MustCompile(`(?i)\b(wife|wives|husbands?)\b`)

var (
	hashtag      = regexp.MustCompile(`(^|\s)#[A-Za-z]`)
	emDash       = regexp.MustCompile(`[ \t]*—[ \t]*`)
	spacedEnDash = regexp.MustCompile(`[ \t]+–[ \t]*`)
)

// findProblems lists every way the lines under the quotes break the voice rules.
// Empty means good to send.
func findProblems(quotes []Quote, lines []string) []string {
	var problems []string
	if len(lines) != len(quotes) {
		problems = append(problems, fmt.Sprintf("there were %d lines for %d quotes", len(lines), len(quotes)))
	}
	calendarCount, callCount := 0, 0
	for _, line := range lines {
		if calendarTalk.MatchString(line) {
			calendarCount++
		}
		if callTalk.MatchString(line) {
			callCount++
		}
	}
	if calendarCount > 1 {
		problems = append(problems, fmt.Sprintf("%d lines mention the day or the time of year, only one may", calendarCount))
	}
	if callCount > 1 {
		problems = append(problems, fmt.Sprintf("%d lines talk about calls or the phones, only one may tie back to the job", callCount))
	}
	for i, line := range lines {
		n := i + 1
		lower := strings.ToLower(line)
		if strings.ContainsAny(line, "—–") {
			problems = append(problems, fmt.Sprintf("line %d uses a dash (— or –)", n))
		}
		for _, phrase := range bannedPhrases {
			if strings.Contains(lower, phrase) {
				problems = append(problems, fmt.Sprintf("line %d says %q", n, phrase))
			}
		}
		if i < len(quotes) && repeatsQuote(quotes[i].Text, line) {
			problems = append(problems, fmt.Sprintf("line %d repeats words from its quote, it should add something new", n))
		}
		if wifeOrHusband.MatchString(line) {
			problems = append(problems, fmt.Sprintf("line %d says wife or husband instead of spouse", n))
		}
		if hashtag.MatchString(line) {
			problems = append(problems, fmt.Sprintf("line %d has a hashtag", n))
		}
		if emojiCount(line) > 1 {
			problems = append(problems, fmt.Sprintf("line %d has more than one emoji", n))
		}
		if strings.Count(line, "!") > 1 {
			problems = append(problems, fmt.Sprintf("line %d has more than one exclamation point", n))
		}
		if strings.Contains(line, "**") || strings.Contains(line, "\"") {
			problems = append(problems, fmt.Sprintf("line %d uses markdown or quote marks, it should be plain talk", n))
		}
		if words := len(strings.Fields(line)); words > maxLineWords {
			problems = append(problems, fmt.Sprintf("line %d is %d words, keep it under 25", n, words))
		}
	}
	return problems
}

// repeatsQuote reports whether the line reuses a run of four words from the quote
// (or the whole quote, if it's shorter), which is how restating it starts.
func repeatsQuote(quote, line string) bool {
	q, l := wordsOf(quote), " "+strings.Join(wordsOf(line), " ")+" "
	size := min(4, len(q))
	for i := 0; i+size <= len(q); i++ {
		if strings.Contains(l, " "+strings.Join(q[i:i+size], " ")+" ") {
			return true
		}
	}
	return false
}

func wordsOf(s string) []string {
	return strings.FieldsFunc(strings.ToLower(strings.ReplaceAll(s, "'", "")), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
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
	return "\nA first try for today came back with these problems. Write all four lines fresh and avoid them:\n- " +
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

// buildDrafts pairs each quote with its line, dashes cleaned up. A missing line just
// leaves the quote on its own.
func buildDrafts(quotes []Quote, lines []string) []Draft {
	drafts := make([]Draft, len(quotes))
	for i, q := range quotes {
		drafts[i] = Draft{Quote: q}
		if i < len(lines) {
			drafts[i].Line = strings.TrimSpace(cleanText(lines[i]))
		}
	}
	return drafts
}
