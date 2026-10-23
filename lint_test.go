package main

import (
	"strings"
	"testing"
)

func goodDrafts() []Draft {
	return []Draft{
		{Label: "every call", Text: "Every call is a door that's already open.\nWalk through it."},
		{Label: "listen first", Text: "Let them talk. The need shows up on its own."},
		{Label: "whole household", Text: "Who else is in that house? Ask. 💪"},
		{Label: "the mission", Text: strings.Repeat("Every policy is a family protected. ", 25)},
	}
}

func TestFindProblemsPassesGoodDrafts(t *testing.T) {
	if problems := findProblems(goodDrafts()); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

func TestFindProblemsCatchesTells(t *testing.T) {
	drafts := goodDrafts()
	drafts[0].Text = "This is a game-changer — trust me! Really! Truly! #grind 🔥💪"
	got := strings.Join(findProblems(drafts), "\n")
	for _, want := range []string{"dash", `"game-changer"`, "hashtag", "emoji", "exclamation"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected a problem mentioning %s, got:\n%s", want, got)
		}
	}
}

func TestFindProblemsWantsSpouse(t *testing.T) {
	for _, text := range []string{"Is your wife covered too?", "Ask about her husband.", "Husbands, wives, kids."} {
		drafts := goodDrafts()
		drafts[1].Text = text
		if got := strings.Join(findProblems(drafts), "\n"); !strings.Contains(got, "spouse") {
			t.Errorf("%q: expected a spouse problem, got:\n%s", text, got)
		}
	}
	drafts := goodDrafts()
	drafts[1].Text = "Is your spouse covered too?"
	if problems := findProblems(drafts); len(problems) != 0 {
		t.Errorf("expected no problems, got %v", problems)
	}
}

func TestFindProblemsCatchesGymWallLines(t *testing.T) {
	drafts := goodDrafts()
	drafts[2].Text = "Get out of your comfort zone. Let that sink in."
	got := strings.Join(findProblems(drafts), "\n")
	for _, want := range []string{`"comfort zone"`, `"let that sink in"`} {
		if !strings.Contains(got, want) {
			t.Errorf("expected a problem mentioning %s, got:\n%s", want, got)
		}
	}
}

func TestFindProblemsAllowsOneCalendarMention(t *testing.T) {
	drafts := goodDrafts()
	drafts[0].Text = "It's Friday. Finish it."
	if problems := findProblems(drafts); len(problems) != 0 {
		t.Fatalf("one calendar mention should pass, got %v", problems)
	}
	drafts[1].Text = "Q4 is on the board."
	got := strings.Join(findProblems(drafts), "\n")
	if !strings.Contains(got, "2 options mention the day or the time of year") {
		t.Errorf("expected a calendar problem, got:\n%s", got)
	}
}

func TestFindProblemsChecksCountAndLength(t *testing.T) {
	drafts := goodDrafts()[:3]
	drafts[1].Text = strings.Repeat("word ", 90)
	got := strings.Join(findProblems(drafts), "\n")
	if !strings.Contains(got, "3 options instead of 4") {
		t.Errorf("expected a count problem, got:\n%s", got)
	}
	if !strings.Contains(got, "option 2 is 90 words") {
		t.Errorf("expected a length problem, got:\n%s", got)
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
