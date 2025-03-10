package findsimilitude

import (
	"fmt"

	"github.com/fabienogli/legigpt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/mistral"
	"github.com/tmc/langchaingo/llms/ollama"
)

func initGPT(cfg legigpt.GPTConfiguration) (llms.Model, error) {
	if cfg.Local != nil {
		return ollama.New(ollama.WithModel(*cfg.Local))

	}
	if cfg.Mistral != nil {
		return mistral.New(mistral.WithAPIKey(cfg.Mistral.ApiKey))
	}
	return nil, fmt.Errorf("when initializing gpt")
}
