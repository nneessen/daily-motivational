package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v5"
)

func subjectFor(day time.Time) string {
	return "Morning texts: " + day.Format("Mon, Jan 2")
}

// formatBody lays the options out plainly so each one is easy to copy on a phone.
func formatBody(drafts []Draft) string {
	var b strings.Builder
	for i, d := range drafts {
		if i > 0 {
			b.WriteString("\n-----\n\n")
		}
		size := "short"
		if i >= shortCount {
			size = "longer"
		}
		fmt.Fprintf(&b, "OPTION %d (%s): %s\n\n%s\n", i+1, size, d.Label, strings.TrimSpace(d.Text))
	}
	return b.String()
}

// sendEmail sends the options as a plain-text email through Mailgun.
func sendEmail(ctx context.Context, subject, body string) error {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	domain := os.Getenv("MAILGUN_DOMAIN")
	to := os.Getenv("TO_ADDRESS")
	if apiKey == "" || domain == "" || to == "" {
		return errors.New("MAILGUN_API_KEY, MAILGUN_DOMAIN and TO_ADDRESS must all be set")
	}

	msg := mailgun.NewMessage(domain, "Morning Texts <morning@"+domain+">", subject, body, to)
	msg.SetTracking(false) // no tracking pixels or rewritten links, just the text
	_, err := mailgun.NewMailgun(apiKey).Send(ctx, msg)
	return err
}
