package httputils

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
)

type ResponseLogger struct {
	client Doer
}

func NewResponseLsogger(client Doer) *ResponseLogger {
	return &ResponseLogger{
		client: client,
	}
}

func (a *ResponseLogger) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	slog.Info("requesting", "url", req.URL, "headers", req.Header)
	resp, err := a.client.Do(ctx, req)
	if err != nil {
		return resp, err
	}
	buf := bytes.Buffer{}
	r := io.TeeReader(resp.Body, &buf)
	data, err := io.ReadAll(r)
	defer func() {
		resp.Body = io.NopCloser(&buf)
	}()
	if err != nil {
		return resp, err
	}
	slog.Info("response", "status", resp.Status, "body", string(data))
	return resp, nil
}
