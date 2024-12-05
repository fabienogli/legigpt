package llmx

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
)

type GPT struct {
	model llms.Model
}

func NewGPT(model llms.Model) *GPT {
	return &GPT{
		model: model,
	}
}

func (g *GPT) Summarize(ctx context.Context, toSummarize string) (string, error) {
	response, err := g.model.GenerateContent(ctx, []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextContent{
					Text: fmt.Sprintf(`
Summarize the following text wrapped between the tag <ACCORD></ACCORD>
In French language with maximum 100 characters:
<ACCORD>%s</ACCORD>
The summary will be in french`, toSummarize),
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
	log.Println(response.Choices)
	return response.Choices[0].Content, nil
}
