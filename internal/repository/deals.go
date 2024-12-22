package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fabienogli/legigpt/internal/domain"
)

type store interface {
	Store(ctx context.Context, data []byte) error
	Get(ctx context.Context) ([]byte, error)
}

type DealRepository struct {
	store store
}

func NewDealRepository(store store) *DealRepository {
	return &DealRepository{
		store: store,
	}
}

func (d *DealRepository) Store(ctx context.Context, searchHistory domain.SearchHistory) error {
	data, err := json.MarshalIndent(searchHistory, "", "  ")
	if err != nil {
		return fmt.Errorf("when marshalling: %w", err)
	}
	err = d.store.Store(ctx, data)
	if err != nil {
		return fmt.Errorf("when storing: %w", err)
	}
	return nil
}

func (d *DealRepository) Get(ctx context.Context) (domain.SearchHistory, error) {
	data, err := d.store.Get(ctx)
	if err != nil {
		return domain.SearchHistory{}, nil
	}
	var searchHistory domain.SearchHistory
	err = json.Unmarshal(data, &searchHistory)
	if err != nil {
		return domain.SearchHistory{}, fmt.Errorf("when unmarshall with data %s: %w", string(data), err)
	}
	return searchHistory, nil
}
