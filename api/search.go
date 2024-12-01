package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/fabienogli/legigpt/httputils"
)

type AuthentifiedClient struct {
	Client httputils.Doer
	URL    string
}

func (a *AuthentifiedClient) Search(ctx context.Context, searchKey string) error {
	seach := map[string]any{
		"recherche": map[string]any{
			"filtres": []any{
				map[string]any{
					"valeurs": []string{
						"LOI",
						"ORDONNANCE",
						"ARRETE",
					},
					"facette": "NATURE",
				},
				map[string]any{
					"dates": map[string]any{
						"start": "2015-01-01",
						"end":   "2018-01-31",
					},
					"facette": "DATE_SIGNATURE",
				},
			},
			"sort":                  "SIGNATURE_DATE_DESC",
			"fromAdvancedRecherche": false,
			"secondSort":            "ID",
			"champs": []any{
				map[string]any{
					"criteres": []any{
						map[string]any{
							"proximite": 2,
							"valeur":    "dispositions",
							"criteres": []any{
								map[string]any{
									"valeur":        "soins",
									"operateur":     "ET",
									"typeRecherche": "UN_DES_MOTS",
								},
								map[string]any{
									"proximite":     "3",
									"valeur":        "fonction publique",
									"operateur":     "ET",
									"typeRecherche": "TOUS_LES_MOTS_DANS_UN_CHAMP",
								},
							},
							"operateur":     "ET",
							"typeRecherche": "UN_DES_MOTS",
						},
					},
					"operateur": "ET",
					"typeChamp": "TITLE",
				},
			},
			"pageSize":       10,
			"operateur":      "ET",
			"typePagination": "DEFAUT",
			"pageNumber":     1,
		},
		"fond": "LODA_DATE",
	}
	payload, err := json.Marshal(seach)
	if err != nil {
		return fmt.Errorf("when json.Marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.URL+"/search", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("whn http.NewRequest: %w", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	resp, err := a.Client.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("when Do: %w", err)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("when readAll: %w", err)
	}
	slog.Info("status", resp.Status, "body", string(buf))
	return nil
}

func (a *AuthentifiedClient) Ping(ctx context.Context) error {
	//get forbidden
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.URL+"/search/ping", nil)
	if err != nil {
		return fmt.Errorf("whn http.NewRequest: %w", err)
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	_, err = a.Client.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("err when ping: %w", err)
	}
	return nil
}
