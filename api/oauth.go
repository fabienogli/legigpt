package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fabienogli/legigpt/httputils"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

type OauthConfig struct {
	ClientSecret string
	ClientID     string
	URL          string
}

type OauthClient struct {
	client       httputils.Doer
	url          string
	clientSecret string
	clientID     string
	Token        string
}

func NewOauthClient(cfg OauthConfig, client httputils.Doer) *OauthClient {
	return &OauthClient{
		client:       client,
		url:          cfg.URL,
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
	}
}

func (a *OauthClient) setToken(ctx context.Context, req *http.Request) error {
	token, err := a.RetrievToken(ctx)
	if err != nil {
		return fmt.Errorf("while retrive token: %w", err)
	}
	a.Token = token.AccessToken
	return nil
}

func (a *OauthClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	if a.Token == "" {
		err := a.setToken(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))

	resp, err := a.client.Do(ctx, req)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (a *OauthClient) RetrievToken(ctx context.Context) (tokenResponse, error) {

	payload := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{a.clientID},
		"client_secret": []string{a.clientSecret},
		"scope":         []string{"openid"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.url, strings.NewReader(payload.Encode()))

	if err != nil {
		return tokenResponse{}, fmt.Errorf("when http.NewRequest: %w", err)
	}

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	resp, err := a.client.Do(ctx, req)

	if err != nil {
		return tokenResponse{}, fmt.Errorf("resp: %#v, err: %w", resp, err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return tokenResponse{}, fmt.Errorf("when io.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return tokenResponse{}, fmt.Errorf("wrong status code(%d), resp=%s", resp.StatusCode, data)
	}
	var respp tokenResponse
	err = json.Unmarshal(data, &respp)
	if err != nil {
		return tokenResponse{}, fmt.Errorf("when json.Unmarshal (resp=%s): %w", data, err)
	}
	return respp, nil
}
