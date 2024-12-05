package usecase

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

type GPT struct {
	model       llms.Model
	summarySize int
}

func NewGPT(model llms.Model, summarySize int) *GPT {
	return &GPT{
		model:       model,
		summarySize: summarySize,
	}
}

func (g *GPT) Summarize(ctx context.Context, toSummarize string) (string, error) {
	response, err := g.model.GenerateContent(ctx, []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{
				llms.TextContent{
					Text: fmt.Sprintf(`
You are a specialized assistant in text summarizing.
Your input is a french text
Your output is a french summary  in less than %d characters. The ouptut is in french.
The summary should be clear, short and cover the main points within the text`, g.summarySize),
				},
			},
		},
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextContent{
					Text: fmt.Sprintf(`[START OF TEXT]
%s
[END OF TEXT]`, toSummarize),
				},
			},
		},
	}, llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}))
	if err != nil {
		return "", fmt.Errorf("when generateContent: %w", err)
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices: %w", err)
	}
	return response.Choices[0].Content, nil
}
