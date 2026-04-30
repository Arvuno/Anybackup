package console

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/signer"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/transport"
)

type Executor struct {
	Signer signer.Signer
	Client transport.Client
}

func (e Executor) Execute(ctx context.Context, remote inputs.RemoteOptions, spec RequestSpec) ([]byte, error) {
	base, err := url.Parse(remote.Endpoint)
	if err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid --endpoint")
	}
	rel, err := url.Parse(spec.Path)
	if err != nil {
		return nil, clierrors.New(clierrors.CodeInvalidArgument, clierrors.ExitInvalidArgs, "invalid request path")
	}

	fullURL := base.ResolveReference(rel)
	req, err := http.NewRequestWithContext(ctx, spec.Method, fullURL.String(), bytes.NewReader(spec.Body))
	if err != nil {
		return nil, clierrors.Wrap(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "failed to build request", err)
	}

	for key, values := range spec.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	if len(spec.Body) > 0 && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if err := e.Signer.Sign(req, signer.Credentials{
		TenantID: remote.TenantID,
		AK:       remote.AK,
		SK:       remote.SK,
	}); err != nil {
		return nil, clierrors.Wrap(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "failed to sign request", err)
	}

	resp, err := e.Client.Do(ctx, transport.Request{
		HTTPRequest: req,
		ReadOnly:    spec.ReadOnly,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, clierrors.Wrap(clierrors.CodeBackendInvalidResp, clierrors.ExitTransport, "failed to read backend response", err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	}
	return nil, clierrors.ClassifyHTTPStatus(resp.StatusCode)
}
