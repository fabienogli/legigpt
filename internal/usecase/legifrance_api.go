package usecase

import (
	"context"
	"fmt"

	"github.com/fabienogli/legigpt/internal/domain"
	"github.com/fabienogli/legigpt/pkg/legifranceapi"
)

type legiApiSearcher interface {
	Search(ctx context.Context, search legifranceapi.Search) (legifranceapi.Response, error)
	Consult(ctx context.Context, search legifranceapi.ConsultRequest) (legifranceapi.AccoPayload, error)
}

type Legifrance struct {
	legiSearcher legiApiSearcher
}

func NewLegifrance(legiSearcher legiApiSearcher) *Legifrance {
	return &Legifrance{
		legiSearcher: legiSearcher,
	}
}

func (d *Legifrance) Search(ctx context.Context, query domain.SearchQuery) (domain.DealResult, error) {
	results, err := d.legiSearcher.Search(ctx, legifranceapi.Search{
		Recherche: legifranceapi.Recherche{
			Filtres: []legifranceapi.Filtre{
				// 	{
				// 		Dates: Dates{
				// 			Start: "2015-01-01",
				// 			End:   "2023-01-31",
				// 		},
				// 		Facette: "DATE_SIGNATURE",
				// 	},
			},
			Sort:                  "SIGNATURE_DATE_DESC",
			FromAdvancedRecherche: false,
			SecondSort:            "ID",
			Champs: []legifranceapi.Champ{
				{
					Criteres: []legifranceapi.Critere{
						{
							// Proximite: 2,
							Valeur: "dispositions",
							Criteres: []legifranceapi.Critere{
								{
									Valeur:        query.Title,
									Operateur:     "ET",
									TypeRecherche: "UN_DES_MOTS",
								},
								// {
								// 	// Proximite:     3,
								// 	Valeur:        "congé",
								// 	Operateur:     "ET",
								// 	TypeRecherche: "UN_DES_MOTS",
								// },
							},
							Operateur:     "ET",
							TypeRecherche: "UN_DES_MOTS",
						},
					},
					Operateur: legifranceapi.OperatorAND,
					TypeChamp: legifranceapi.FieldAll,
				},
			},
			PageSize:       query.LimitSize,
			Operateur:      legifranceapi.OperatorAND,
			TypePagination: legifranceapi.PaginationDefault,
			PageNumber:     query.PageNumber,
		},
		Fond: legifranceapi.FondACCO,
	})
	if err != nil {
		return domain.DealResult{}, fmt.Errorf("when search: %w", err)
	}

	var accords []domain.Accord
	for _, result := range results.Results {
		for _, title := range result.Titles {
			accords = append(accords, domain.Accord{
				ID:    title.ID,
				CID:   title.CID,
				Title: title.Title,
			})
		}
	}
	return domain.DealResult{
		Accords: accords,
		Total:   results.TotalResultNumber,
	}, nil
}

func (d *Legifrance) GetContent(ctx context.Context, accord domain.Accord) (domain.Accord, error) {
	results, err := d.legiSearcher.Consult(ctx, legifranceapi.ConsultRequest{
		ID: accord.ID,
	})
	if err != nil {
		return accord, fmt.Errorf("when consult: %w", err)
	}
	accord.Texte = results.Acco.Attachment.Content
	return accord, nil

}
