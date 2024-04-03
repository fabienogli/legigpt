package httputils

import (
	"context"
	"net/http"
)

type Doer interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
