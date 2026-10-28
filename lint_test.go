package main

import (
	"strings"
	"testing"
)

func testQuotes() []Quote {
	return []Quote{
		{Text: "Pressure is a privilege.", Author: "Billie Jean King"},
		{Text: "Never die easy.", Author: "Walter Payton"},
		{Text: "The best way out is always through.", Author: "Robert Frost"},
		{Text: "Well done is better than well said.", Author: "Benjamin Franklin"},
	}
}

func goodLines() []string {
	return []string{
		"The hard call today is the one that pays. Take it.",
		"LET'S WORK. 💪",
		"",
		"Somebody's whole family gets covered because you asked one more question.",
	}
}

func TestFindProblemsPassesGoodLines(t *testing.T) {
	if problems := findProblems(testQuotes(), goodLines()); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

func TestFindProblemsCatchesTells(t *testing.T) {
	lines := goodLines()
	lines[0] = "This is a game-changer — trust me! Really! #grind 🔥💪"
	got := strings.Join(findProblems(testQuotes(), lines), "\n")
	for _, want := range []string{"dash", `"game-changer"`, "hashtag", "emoji", "exclamation"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected a problem mentioning %s, got:\n%s", want, got)
		}
	}
}

func TestFindProblemsWantsSpouse(t *testing.T) {
	for _, line := range []string{"Is your wife covered too?", "Ask about her husband.", "Husbands, wives, kids."} {
		lines := goodLines()
		lines[1] = line
		if got := strings.Join(findProblems(testQuotes(), lines), "\n"); !strings.Contains(got, "spouse") {
			t.Errorf("%q: expected a spouse problem, got:\n%s", line, got)
		}
	}
	lines := goodLines()
	lines[1] = "Is your spouse covered too?"
	if problems := findProblems(testQuotes(), lines); len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}
}

func TestFindProblemsCatchesGymWallLines(t *testing.T) {
	lines := goodLines()
	lines[2] = "Get out of your comfort zone. Let that sink in."
	got := strings.Join(findProblems(testQuotes(), lines), "\n")
	for _, want := range []string{`"comfort zone"`, `"let that sink in"`} {
		if !strings.Contains(got, want) {
			t.Errorf("expected a problem mentioning %s, got:\n%s", want, got)
		}
	}
}

func TestFindProblemsCatchesRestatingTheQuote(t *testing.T) {
	lines := goodLines()
	lines[0] = "Remember, pressure is a privilege, so go get it."
	lines[1] = "Never die easy out there."
	got := strings.Join(findProblems(testQuotes(), lines), "\n")
	for _, want := range []string{"line 1 repeats", "line 2 repeats"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q, got:\n%s", want, got)
		}
	}
	lines = goodLines()
	lines[1] = "Easy days don't build anything."
	if problems := findProblems(testQuotes(), lines); len(problems) != 0 {
		t.Errorf("sharing one word with the quote is fine, got %v", problems)
	}
}

func TestFindProblemsCatchesExtraQuotes(t *testing.T) {
	lines := goodLines()
	lines[0] = `Like my old coach said, "keep going."`
	if got := strings.Join(findProblems(testQuotes(), lines), "\n"); !strings.Contains(got, "quote marks") {
		t.Errorf("expected a quote marks problem, got:\n%s", got)
	}
}

func TestFindProblemsAllowsOneCalendarMention(t *testing.T) {
	lines := goodLines()
	lines[0] = "It's Friday. Finish it."
	if problems := findProblems(testQuotes(), lines); len(problems) != 0 {
		t.Fatalf("one calendar mention should pass, got %v", problems)
	}
	lines[1] = "Q4 is on the board."
	got := strings.Join(findProblems(testQuotes(), lines), "\n")
	if !strings.Contains(got, "2 lines mention the day or the time of year") {
		t.Errorf("expected a calendar problem, got:\n%s", got)
	}
}

func TestFindProblemsAllowsOneLineAboutCalls(t *testing.T) {
	lines := goodLines()
	lines[0] = "Same energy on call forty as call one."
	if problems := findProblems(testQuotes(), lines); len(problems) != 0 {
		t.Fatalf("one line about calls should pass, got %v", problems)
	}
	lines[3] = "Every caller has a family behind them."
	got := strings.Join(findProblems(testQuotes(), lines), "\n")
	if !strings.Contains(got, "2 lines talk about calls") {
		t.Errorf("expected a calls problem, got:\n%s", got)
	}
}

func TestFindProblemsChecksCountAndLength(t *testing.T) {
	lines := goodLines()[:3]
	lines[1] = strings.Repeat("word ", 40)
	got := strings.Join(findProblems(testQuotes(), lines), "\n")
	if !strings.Contains(got, "3 lines for 4 quotes") {
		t.Errorf("expected a count problem, got:\n%s", got)
	}
	if !strings.Contains(got, "line 2 is 40 words") {
		t.Errorf("expected a length problem, got:\n%s", got)
	}
}

func TestBuildDraftsHandlesMissingLines(t *testing.T) {
	drafts := buildDrafts(testQuotes(), []string{"Go — now."})
	if len(drafts) != 4 {
		t.Fatalf("got %d drafts, want 4", len(drafts))
	}
	if drafts[0].Line != "Go... now." {
		t.Errorf("dash not cleaned: %q", drafts[0].Line)
	}
	if drafts[3].Line != "" || drafts[3].Text() != "\"Well done is better than well said.\"\n- Benjamin Franklin" {
		t.Errorf("a missing line should leave the quote alone, got %q", drafts[3].Text())
	}
}

func TestCleanText(t *testing.T) {
	tests := map[string]string{
		"Every call counts — all of them.": "Every call counts... all of them.",
		"Every call counts—all of them.":   "Every call counts... all of them.",
		"Take the call –\nevery one.":      "Take the call...\nevery one.",
		"Keep it to 2–5 lines.":            "Keep it to 2-5 lines.",
	}
	for in, want := range tests {
		if got := cleanText(in); got != want {
			t.Errorf("cleanText(%q) = %q, want %q", in, got, want)
		}
	}
}
