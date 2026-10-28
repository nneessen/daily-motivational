package main

import (
	"cmp"
	_ "embed"
	"fmt"
	"hash/fnv"
	"slices"
	"strings"
	"time"
)

// quotesFile is the quote bank. Edit quotes.txt to add or remove quotes.
//
//go:embed quotes.txt
var quotesFile string

// Quote is one real quote from the bank, printed exactly as written there.
type Quote struct {
	Text   string
	Author string
}

// parseQuotes reads "quote | who said it" lines, skipping blanks and # comments. Lines
// that don't fit come back in bad, so a typo made on github.com never costs a morning.
func parseQuotes(raw string) (quotes []Quote, bad []string) {
	for n, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		text, author, ok := strings.Cut(line, " | ")
		text, author = strings.TrimSpace(text), strings.TrimSpace(author)
		if !ok || text == "" || author == "" || strings.Contains(author, "|") {
			bad = append(bad, fmt.Sprintf("line %d: %q", n+1, line))
			continue
		}
		quotes = append(quotes, Quote{Text: text, Author: author})
	}
	return quotes, bad
}

// rotationStart is a Monday. Weekdays are counted from here to pick each day's quotes.
var rotationStart = time.Date(2026, time.October, 26, 0, 0, 0, 0, time.UTC)

// weekdayNumber counts the weekdays from rotationStart to day. A Saturday or Sunday
// shares the next Monday's number.
func weekdayNumber(day time.Time) int {
	date := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	days := int(date.Sub(rotationStart).Hours() / 24)
	return days/7*5 + min(days%7, 5)
}

// quotesFor picks the day's quotes. The bank is put in a fixed order by a hash of each
// quote and each weekday takes the next n, so no quote comes back until the whole bank
// has been used. Because the order comes from the quotes themselves, adding or removing
// one only moves that one, and at most a quote or two from recent mornings comes back
// early. Rerunning on the same day gives the same quotes.
func quotesFor(day time.Time, bank []Quote, n int) []Quote {
	order := slices.Clone(bank)
	slices.SortFunc(order, func(a, b Quote) int { return cmp.Compare(rotationKey(a), rotationKey(b)) })
	start := weekdayNumber(day) * n
	picked := make([]Quote, n)
	for i := range picked {
		picked[i] = order[((start+i)%len(order)+len(order))%len(order)]
	}
	return picked
}

// rotationKey is where a quote sits in the rotation.
func rotationKey(q Quote) uint64 {
	h := fnv.New64a()
	h.Write([]byte(q.Text))
	return h.Sum64()
}
