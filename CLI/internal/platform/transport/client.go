package transport

import (
	"context"
	"errors"
	"net/http"

	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type DoerFunc func(req *http.Request) (*http.Response, error)

func (f DoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type Request struct {
	HTTPRequest *http.Request
	ReadOnly    bool
}

type Client struct {
	Doer Doer
}

func (c Client) Do(ctx context.Context, req Request) (*http.Response, error) {
	attempts := 1
	if req.ReadOnly {
		attempts = 2
	}

	var lastErr error
	for i := 0; i < attempts; i++ {
		resp, err := c.Doer.Do(req.HTTPRequest.WithContext(ctx))
		if err == nil {
			return resp, nil
		}
		lastErr = err
	}

	if errors.Is(lastErr, context.DeadlineExceeded) {
		return nil, clierrors.Wrap(clierrors.CodeTransportTimeout, clierrors.ExitTransport, "request timed out", lastErr)
	}
	return nil, clierrors.Wrap(clierrors.CodeTransportConnection, clierrors.ExitTransport, "request failed", lastErr)
}
