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

type DealLooker struct {
	legiSearcher legiSearcher
}

func NewDealLooker(legiSearcher legiSearcher) *DealLooker {
	return &DealLooker{
		legiSearcher: legiSearcher,
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
