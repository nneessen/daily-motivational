package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

// Draft is one text option: a quote from the bank and the line Nick adds under it.
type Draft struct {
	Quote Quote
	Line  string
}

// Text is the message as Nick would send it. The quote comes straight from the bank,
// so Claude can never misquote anyone.
func (d Draft) Text() string {
	text := fmt.Sprintf("\"%s\"\n- %s", d.Quote.Text, d.Quote.Author)
	if line := strings.TrimSpace(d.Line); line != "" {
		text += "\n\n" + line
	}
	return text
}

// linesSchema is the JSON shape Claude has to answer in, so there's never a
// "Here are your messages!" preamble to strip off.
var linesSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"lines": map[string]any{
			"type":  "array",
			"items": map[string]any{"type": "string"},
		},
	},
	"required":             []string{"lines"},
	"additionalProperties": false,
}

// generate asks Claude for the line under each quote, one per quote, in order.
func generate(ctx context.Context, client *anthropic.Client, prompt string) ([]string, error) {
	resp, err := client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		Model:     "claude-opus-5-5",
		MaxTokens: 16000,
		System:    []anthropic.BetaTextBlockParam{{Text: voiceGuide}},
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(prompt)),
		},
		OutputConfig: anthropic.BetaOutputConfigParam{
			Effort: anthropic.BetaOutputConfigEffortHigh,
			Format: anthropic.BetaJSONOutputFormatParam{Schema: linesSchema},
		},
		// If a safety check ever declines the request, the API reruns it on a backup model.
		Fallbacks: anthropic.BetaFallbacksParamOfDefault(),
		Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaServerSideFallback2026_07_01},
	})
	if err != nil {
		return nil, fmt.Errorf("claude request: %w", err)
	}

	switch resp.StopReason {
	case anthropic.BetaStopReasonRefusal:
		return nil, fmt.Errorf("claude declined the request: %s", resp.StopDetails.Explanation)
	case anthropic.BetaStopReasonMaxTokens:
		return nil, errors.New("claude hit the output limit before finishing")
	}

	for _, block := range resp.Content {
		if text, ok := block.AsAny().(anthropic.BetaTextBlock); ok {
			var out struct {
				Lines []string `json:"lines"`
			}
			if err := json.Unmarshal([]byte(text.Text), &out); err != nil {
				return nil, fmt.Errorf("parse claude output: %w", err)
			}
			return out.Lines, nil
		}
	}
	return nil, errors.New("claude response had no text block")
}
