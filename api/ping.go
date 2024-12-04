package api

import (
	"context"
	"fmt"
	"net/http"
)

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
