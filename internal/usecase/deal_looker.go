package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/fabienogli/legigpt/internal/domain"
)

type legiSearcher interface {
	Search(ctx context.Context, query domain.SearchQuery) (domain.DealResult, error)
	GetContent(context.Context, domain.Accord) (domain.Accord, error)
}

type gpt interface {
	Summarize(ctx context.Context, toSummarize string) (string, error)
	FindSimilitude(ctx context.Context, knowledge []string, delimiter, searchSimilarity string) (string, error)
}

type DealLooker struct {
	legiSearcher legiSearcher
	gpt          gpt
}

func NewDealLooker(legiSearcher legiSearcher, gpt gpt) *DealLooker {
	return &DealLooker{
		legiSearcher: legiSearcher,
		gpt:          gpt,
	}
}

func (d *DealLooker) Search(ctx context.Context, query domain.SearchQuery) (domain.DealResult, error) {
	accords, err := d.legiSearcher.Search(ctx, query)
	if err != nil {
		return domain.DealResult{}, fmt.Errorf("when dealLooker.Search: %w", err)
	}
	log.Println("total: ", accords.Total)

	for i, accord := range accords.Accords {
		content, err := d.legiSearcher.GetContent(ctx, accord)
		if err != nil {
			return domain.DealResult{}, fmt.Errorf("when dealLooker.GetContent: %w", err)
		}
		accords.Accords[i] = content
	}
	return accords, nil
}

func (d *DealLooker) Rag(ctx context.Context, deals []domain.Accord, keyword string) (string, error) {
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
