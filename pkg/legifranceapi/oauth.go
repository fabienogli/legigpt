package legifranceapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/fabienogli/legigpt/httputils"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	// ExpiredInSecond int    `json:"expires_in"`
	//TODO add time
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
	filestore    cache
}

type cache interface {
	Store(ctx context.Context, data []byte) error
	Get(ctx context.Context) ([]byte, error)
}

func NewOauthClient(cfg OauthConfig, client httputils.Doer, filestore cache) *OauthClient {
	return &OauthClient{
		client:       client,
		url:          cfg.URL,
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		filestore:    filestore,
	}
}

func (a *OauthClient) getTokenFromFile(ctx context.Context) (TokenResponse, error) {
	var resp TokenResponse
	data, err := a.filestore.Get(ctx)
	if err != nil {
		return resp, fmt.Errorf("when a.filestore.Get: %w", err)
	}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, fmt.Errorf("when json.Unmarshal: %w", err)
	}
	return resp, nil
}

func (a *OauthClient) saveTokenFromFile(ctx context.Context, token TokenResponse) error {
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("when json.marshal: %w", err)
	}
	err = a.filestore.Store(ctx, data)
	if err != nil {
		return fmt.Errorf("when a.filestore.Store: %w", err)
	}
	return nil
}

func (a *OauthClient) SetToken(ctx context.Context, req *http.Request) error {
	token, err := a.getTokenFromFile(ctx)
	if err != nil {
		token, err := a.retrievToken(ctx)
		if err != nil {
			return fmt.Errorf("while retrive token: %w", err)
		}
		err = a.saveTokenFromFile(ctx, token)
		if err != nil {
			log.Println("err when saving token: %v", err)
		}
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	return nil
}

// Do will retrieve the token if not set
func (a *OauthClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	a.SetToken(ctx, req)
	resp, err := a.client.Do(ctx, req)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (a *OauthClient) retrievToken(ctx context.Context) (TokenResponse, error) {

	payload := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{a.clientID},
		"client_secret": []string{a.clientSecret},
		"scope":         []string{"openid"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.url, strings.NewReader(payload.Encode()))

	if err != nil {
		return TokenResponse{}, fmt.Errorf("when http.NewRequest: %w", err)
	}

	req.Header.Add("Content-type", "application/x-www-form-urlencoded")
	resp, err := a.client.Do(ctx, req)

	if err != nil {
		return TokenResponse{}, fmt.Errorf("resp: %#v, err: %w", resp, err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("when io.ReadAll: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return TokenResponse{}, fmt.Errorf("wrong status code(%d), resp=%s", resp.StatusCode, data)
	}
	var respp TokenResponse
	err = json.Unmarshal(data, &respp)
	if err != nil {
		return TokenResponse{}, fmt.Errorf("when json.Unmarshal (resp=%s): %w", data, err)
	}
	return respp, nil
}
