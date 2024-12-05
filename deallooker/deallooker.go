package deallooker

import (
	"context"
	"fmt"

	"github.com/fabienogli/legigpt/api"
)

type legiApiSearcher interface {
	Search(ctx context.Context, search api.Search) (api.Response, error)
	Consult(ctx context.Context, search api.ConsultRequest) (api.AccoPayload, error)
}

type DealLooker struct {
	legiSearcher legiApiSearcher
}

type SearchQuery struct {
	Title string
	//PageSize Max=100
	//TODO: limit and offset
	PageSize   int
	PageNumber int
}

func NewDealLooker(legiSearcher legiApiSearcher) *DealLooker {
	return &DealLooker{
		legiSearcher: legiSearcher,
	}
}

func (d *DealLooker) Search(ctx context.Context, query SearchQuery) (AccordsWrapped, error) {
	results, err := d.legiSearcher.Search(ctx, api.Search{
		Recherche: api.Recherche{
			Filtres: []api.Filtre{
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
			Champs: []api.Champ{
				{
					Criteres: []api.Critere{
						{
							// Proximite: 2,
							Valeur: "dispositions",
							Criteres: []api.Critere{
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
					Operateur: api.OperatorAND,
					TypeChamp: api.FieldAll,
				},
			},
			PageSize:       query.PageSize,
			Operateur:      api.OperatorAND,
			TypePagination: api.PaginationDefault,
			PageNumber:     query.PageNumber,
		},
		Fond: api.FondACCO,
	})
	if err != nil {
		return AccordsWrapped{}, fmt.Errorf("when search: %w", err)
	}

	var accords []Accord
	for _, result := range results.Results {
		for _, title := range result.Titles {
			accords = append(accords, Accord{
				ID:    title.ID,
				CID:   title.CID,
				Title: title.Title,
			})
		}
	}
	return AccordsWrapped{
		Accords: accords,
		Total:   results.TotalResultNumber,
	}, nil
}

func (d *DealLooker) GetContent(ctx context.Context, id string) (Content, error) {
	results, err := d.legiSearcher.Consult(ctx, api.ConsultRequest{
		ID: id,
	})
	if err != nil {
		return Content{}, fmt.Errorf("when consult: %w", err)
	}
	return Content{
		Texte: results.Acco.Attachment.Content,
	}, nil

}
