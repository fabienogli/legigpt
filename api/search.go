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
		"fond": "ACCO",
		"recherche": map[string]any{
			"operateur":      "ET",
			"pageSize":       10,
			"sort":           "SIGNATURE_DATE_DESC",
			"typePagination": "DEFAULT",
			"pageNumber":     1,
			"champs": []map[string]any{
				{
					"operateur": "ET",
					"criteres": []map[string]any{
						{
							"operateur":     "ET",
							"valeur":        searchKey,
							"typeRecherche": "UN_DES_MOTS",
						},
					},
				},
			},
			"typeChamp": "TITLE",
		},
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
