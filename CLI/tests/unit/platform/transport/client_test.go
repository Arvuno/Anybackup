package transport_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/transport"
)

func TestDo_RetriesReadOnlyTimeoutButNotWrite(t *testing.T) {
	attempts := 0
	client := transport.Client{
		Doer: transport.DoerFunc(func(req *http.Request) (*http.Response, error) {
			attempts++
			return nil, context.DeadlineExceeded
		}),
	}

	_, err := client.Do(context.Background(), transport.Request{
		HTTPRequest: mustRequest(t, http.MethodGet),
		ReadOnly:    true,
	})
	if err == nil {
		t.Fatal("Do() error = nil, want timeout")
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}

	attempts = 0
	_, _ = client.Do(context.Background(), transport.Request{
		HTTPRequest: mustRequest(t, http.MethodPost),
		ReadOnly:    false,
	})
	if attempts != 1 {
		t.Fatalf("write attempts = %d, want 1", attempts)
	}
}

func mustRequest(t *testing.T, method string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, "https://127.0.0.1:9600/x", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
