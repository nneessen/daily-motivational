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
	flag.Parse()

	denver, err := time.LoadLocation("America/Denver")
	if err != nil {
		log.Fatalf("load timezone: %v", err)
	}
	today := time.Now().In(denver)

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

// generateChecked asks for today's drafts. If they break the voice rules, it asks once
// more with the problems spelled out. A second miss still goes out, dashes cleaned up,
// so a bad draft never costs a morning.
func generateChecked(ctx context.Context, client *anthropic.Client, today time.Time) ([]Draft, error) {
	prompt := buildPrompt(today, pickAngles(draftCount), readFocus("focus.txt"))
	drafts, err := generate(ctx, client, prompt)
	if err != nil {
		return nil, err
	}
	problems := findProblems(drafts)
	if len(problems) == 0 {
		return drafts, nil
	}

	log.Printf("first draft had problems, retrying: %v", problems)
	retry, err := generate(ctx, client, prompt+retryNote(problems))
	if err != nil || len(retry) == 0 {
		log.Printf("retry failed (%v), sending the first draft", err)
		return cleanDrafts(drafts), nil
	}
	if problems := findProblems(retry); len(problems) > 0 {
		log.Printf("retry still had problems, sending anyway: %v", problems)
	}
	return cleanDrafts(retry), nil
}
