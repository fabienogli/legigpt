package usecase

import (
	"context"
	"fmt"

	"github.com/fabienogli/legigpt/internal/domain"
)

type gpt interface {
	Summarize(ctx context.Context, toSummarize string) (string, error)
	FindSimilitude(ctx context.Context, knowledge []string, delimiter, searchSimilarity string) (string, error)
}

type DealLLM struct {
	gpt gpt
}

func NewDealLLM(gpt gpt) *DealLLM {
	return &DealLLM{
		gpt: gpt,
	}
}
func (d *DealLLM) Rag(ctx context.Context, deals []domain.Accord, keyword string) (string, error) {
	//joining the deals so you can wrap it in a prompt
	delimiter := "________________________________________"
	textJoin := make([]string, len(deals))
	for i, deal := range deals {
		textJoin[i] = deal.Texte
	}

	bestDeal, err := d.gpt.FindSimilitude(ctx, textJoin, delimiter, keyword)

	if err != nil {
		return "", fmt.Errorf("when finding similiraties: %w", err)
	}

	return bestDeal, nil
}
