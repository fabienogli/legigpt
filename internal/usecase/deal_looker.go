package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/fabienogli/legigpt/internal/domain"
)

type legiSearcher interface {
	Search(ctx context.Context, query domain.SearchQuery) (domain.AccordsWrapped, error)
	GetContent(ctx context.Context, id string) (domain.Content, error)
}

type gpt interface {
	Summarize(ctx context.Context, toSummarize string) (string, error)
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

func (d *DealLooker) Search(ctx context.Context, query domain.SearchQuery) error {
	// accords, err := d.legiSearcher.Search(ctx, domain.SearchQuery{
	// 	Title:      "congé",
	// 	PageSize:   1,
	// 	PageNumber: 1,
	// })
	accords, err := d.legiSearcher.Search(ctx, query)
	if err != nil {
		return fmt.Errorf("when dealLooker.Search: %w", err)
	}
	log.Println("total: ", accords.Total)

	var contents []domain.Content
	for _, accord := range accords.Accords {
		content, err := d.legiSearcher.GetContent(ctx, accord.ID)
		if err != nil {
			return fmt.Errorf("when dealLooker.GetContent: %w", err)
		}
		summary, err := d.gpt.Summarize(ctx, content.Texte)
		if err != nil {
			return fmt.Errorf("when ollamaTest: %w", err)
		}
		log.Println("Summary: ", summary)
		contents = append(contents, content)
	}
	return nil
}
