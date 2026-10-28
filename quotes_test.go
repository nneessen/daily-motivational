package main

import (
	"slices"
	"strings"
	"testing"
	"time"
)

func TestQuoteBankParses(t *testing.T) {
	bank, bad := parseQuotes(quotesFile)
	if len(bad) > 0 {
		t.Errorf("quotes.txt lines should be \"quote | who said it\": %v", bad)
	}
	if len(bank) < 120 {
		t.Fatalf("only %d quotes in the bank, want at least 120 so a quote doesn't come back for six weeks", len(bank))
	}
	seen := map[string]bool{}
	for _, q := range bank {
		if seen[q.Text] {
			t.Errorf("quote listed twice: %q", q.Text)
		}
		seen[q.Text] = true
		if strings.ContainsAny(q.Text+q.Author, "“”—–") {
			t.Errorf("use plain quotes and hyphens, Nick's phone shows the fancy ones oddly: %q", q.Text)
		}
	}
}

func TestParseQuotesRejectsBadLines(t *testing.T) {
	got, bad := parseQuotes("# comment\n\nA line with no author\nWell done is better than well said. | Benjamin Franklin\n")
	if len(got) != 1 || got[0].Author != "Benjamin Franklin" {
		t.Errorf("parseQuotes kept %v, want just the Franklin line", got)
	}
	if len(bad) != 1 || !strings.Contains(bad[0], "line 3") {
		t.Errorf("parseQuotes bad = %v, want line 3 flagged", bad)
	}
}

func TestWeekdayNumber(t *testing.T) {
	tests := map[string]int{
		"2026-10-26": 0, // Monday
		"2026-10-30": 4, // Friday
		"2026-10-31": 5, // Saturday shares the next Monday's number
		"2026-11-02": 5, // Monday
		"2026-12-25": 44,
	}
	for date, want := range tests {
		day, _ := time.Parse("2006-01-02", date)
		if got := weekdayNumber(day); got != want {
			t.Errorf("weekdayNumber(%s) = %d, want %d", date, got, want)
		}
	}
}

func TestQuotesDontRepeatUntilTheBankRunsOut(t *testing.T) {
	bank, _ := parseQuotes(quotesFile)
	weekdays := len(bank) / draftCount
	seen := map[string]string{}
	day := time.Date(2026, 10, 28, 5, 17, 0, 0, time.UTC)
	for count := 0; count < weekdays; day = day.AddDate(0, 0, 1) {
		if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
			continue
		}
		for _, q := range quotesFor(day, bank, draftCount) {
			if first, ok := seen[q.Text]; ok {
				t.Fatalf("%q came back on %s, first used %s", q.Text, day.Format("Jan 2"), first)
			}
			seen[q.Text] = day.Format("Jan 2")
		}
		count++
	}
}

func TestQuotesForIsSteadyWithinADay(t *testing.T) {
	bank, _ := parseQuotes(quotesFile)
	morning := time.Date(2026, 10, 28, 5, 17, 0, 0, time.UTC)
	rerun := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	a, b := quotesFor(morning, bank, draftCount), quotesFor(rerun, bank, draftCount)
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("a rerun the same day picked different quotes: %v vs %v", a, b)
		}
	}
}

// weekdaysFrom lists count weekdays starting at day.
func weekdaysFrom(day time.Time, count int) []time.Time {
	var days []time.Time
	for ; len(days) < count; day = day.AddDate(0, 0, 1) {
		if day.Weekday() != time.Saturday && day.Weekday() != time.Sunday {
			days = append(days, day)
		}
	}
	return days
}

func TestAddingAQuoteBarelyChangesTheRotation(t *testing.T) {
	bank, _ := parseQuotes(quotesFile)
	days := weekdaysFrom(time.Date(2026, 10, 28, 5, 17, 0, 0, time.UTC), 30)
	sent := map[Quote]bool{}
	for _, day := range days[:10] {
		for _, q := range quotesFor(day, bank, draftCount) {
			sent[q] = true
		}
	}
	// Nick adds a quote on github.com, then the next 20 mornings run.
	edited := append(slices.Clone(bank), Quote{Text: "A brand new quote nobody has seen.", Author: "Someone"})
	repeats := 0
	for _, day := range days[10:] {
		for _, q := range quotesFor(day, edited, draftCount) {
			if sent[q] {
				repeats++
			}
		}
	}
	if repeats > 1 {
		t.Errorf("adding one quote brought back %d quotes from the last 10 mornings, want at most 1", repeats)
	}
}
