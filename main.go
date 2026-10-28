// Command daily-motivational emails Nick a few options for his morning text to the team.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
	_ "time/tzdata" // embed the tz database so Denver time never depends on the machine

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "print the email instead of sending it")
	date := flag.String("date", "", "write the texts for another day, like 2026-11-02 (handy with -dry-run)")
	flag.Parse()

	denver, err := time.LoadLocation("America/Denver")
	if err != nil {
		log.Fatalf("load timezone: %v", err)
	}
	today := time.Now().In(denver)
	if *date != "" {
		if today, err = time.ParseInLocation("2006-01-02", *date, denver); err != nil {
			log.Fatalf("-date should look like 2026-11-02: %v", err)
		}
	}

	ctx := context.Background()
	client := anthropic.NewClient()
	drafts, err := generateChecked(ctx, &client, today)
	if err != nil {
		log.Fatal(err)
	}

	subject := subjectFor(today)
	body := formatBody(drafts)
	if *dryRun {
		fmt.Printf("Subject: %s\n\n%s", subject, body)
		return
	}
	if err := sendEmail(ctx, subject, body); err != nil {
		log.Fatalf("send email: %v", err)
	}
	log.Printf("sent %q", subject)
}

// generateChecked picks today's quotes and asks for the line under each one. If the
// lines break the voice rules, it asks once more with the problems spelled out. A second
// miss still goes out, dashes cleaned up, so a bad line never costs a morning.
func generateChecked(ctx context.Context, client *anthropic.Client, today time.Time) ([]Draft, error) {
	bank, bad := parseQuotes(quotesFile)
	if len(bad) > 0 {
		log.Printf("skipping quotes.txt lines that aren't \"quote | who said it\": %v", bad)
	}
	if len(bank) < draftCount {
		return nil, fmt.Errorf("quotes.txt has %d usable quotes, it needs at least %d", len(bank), draftCount)
	}
	quotes := quotesFor(today, bank, draftCount)
	prompt := buildPrompt(today, quotes, readFocus("focus.txt"))
	lines, err := generate(ctx, client, prompt)
	if err != nil {
		return nil, err
	}
	problems := findProblems(quotes, lines)
	if len(problems) == 0 {
		return buildDrafts(quotes, lines), nil
	}

	log.Printf("first try had problems, retrying: %v", problems)
	retry, err := generate(ctx, client, prompt+retryNote(problems))
	if err != nil || len(retry) == 0 {
		log.Printf("retry failed (%v), sending the first try", err)
		return buildDrafts(quotes, lines), nil
	}
	if problems := findProblems(quotes, retry); len(problems) > 0 {
		log.Printf("retry still had problems, sending anyway: %v", problems)
	}
	return buildDrafts(quotes, retry), nil
}
