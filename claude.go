package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// Draft is one text option for the team.
type Draft struct {
	Label string `json:"label"`
	Text  string `json:"text"`
}

// draftSchema is the JSON shape Claude has to answer in, so there's never a
// "Here are your messages!" preamble to strip off.
var draftSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"messages": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"label": map[string]any{"type": "string"},
					"text":  map[string]any{"type": "string"},
				},
				"required":             []string{"label", "text"},
				"additionalProperties": false,
			},
		},
	},
	"required":             []string{"messages"},
	"additionalProperties": false,
}

func generate(ctx context.Context, client *anthropic.Client, prompt string) ([]Draft, error) {
	resp, err := client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		Model:     "claude-opus-5-5",
		MaxTokens: 16000,
		System:    []anthropic.BetaTextBlockParam{{Text: voiceGuide}},
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(prompt)),
		},
		OutputConfig: anthropic.BetaOutputConfigParam{
			Effort: anthropic.BetaOutputConfigEffortHigh,
			Format: anthropic.BetaJSONOutputFormatParam{Schema: draftSchema},
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
				Messages []Draft `json:"messages"`
			}
			if err := json.Unmarshal([]byte(text.Text), &out); err != nil {
				return nil, fmt.Errorf("parse claude output: %w", err)
			}
			return out.Messages, nil
		}
	}
	return nil, errors.New("claude response had no text block")
}
