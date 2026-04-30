package console_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/business/inputs"
	clierrors "github.com/anybackup-ai/Anybackup/CLI/internal/errors"
	"github.com/anybackup-ai/Anybackup/CLI/internal/external/console"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/signer"
	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/transport"
)

func TestExecute_ReturnsRawBodyAndClassifiesStatuses(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/unauthorized":
			http.Error(w, "no auth", http.StatusUnauthorized)
		default:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("{\"ok\":true}"))
		}
	}))
	defer srv.Close()

	exec := console.Executor{
		Signer: signer.Signer{},
		Client: transport.Client{Doer: srv.Client()},
	}

	body, err := exec.Execute(context.Background(), inputs.RemoteOptions{
		TenantID:      "tenant-a",
		Endpoint:      srv.URL,
		AK:            "ak",
		SK:            "sk",
		TargetVersion: "9.0.9.0",
	}, console.RequestSpec{
		Method:   http.MethodGet,
		Path:     "/api/sla/v1/templates",
		ReadOnly: true,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if string(body) != "{\"ok\":true}" {
		t.Fatalf("body = %q", string(body))
	}

	_, err = exec.Execute(context.Background(), inputs.RemoteOptions{
		TenantID: "tenant-a", Endpoint: srv.URL, AK: "ak", SK: "sk", TargetVersion: "9.0.9.0",
	}, console.RequestSpec{
		Method:   http.MethodGet,
		Path:     "/unauthorized",
		ReadOnly: true,
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want auth error")
	}
	if clierrors.ExitCodeOf(err) != clierrors.ExitAuth {
		t.Fatalf("ExitCodeOf(err) = %d, want %d", clierrors.ExitCodeOf(err), clierrors.ExitAuth)
	}
}

func TestExecute_SendsSDKStyleSignatureHeaders(t *testing.T) {
	var gotTenant string
	var gotAccessKey string
	var gotTimestamp string
	var gotSignature map[string]string
	var gotBody []byte
	var decodeErr error
	var gotSignatureRaw string
	var bodyReadErr error

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTenant = r.Header.Get("X-Tenant-Id")
		gotAccessKey = r.Header.Get("X-Access-Key")
		gotTimestamp = r.Header.Get("X-Timestamp")
		gotSignatureRaw = r.Header.Get("Signature")
		gotBody, bodyReadErr = io.ReadAll(r.Body)
		gotSignature = map[string]string{}
		decodeErr = json.Unmarshal([]byte(gotSignatureRaw), &gotSignature)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	exec := console.Executor{
		Signer: signer.Signer{},
		Client: transport.Client{Doer: srv.Client()},
	}

	body, err := exec.Execute(context.Background(), inputs.RemoteOptions{
		TenantID:      "tenant-a",
		Endpoint:      srv.URL,
		AK:            "ak",
		SK:            "sk",
		TargetVersion: "9.0.9.0",
	}, console.RequestSpec{
		Method:   http.MethodPost,
		Path:     "/api/file%20recovery?scope=a%2Fb",
		Body:     []byte(`{"name":"daily"}`),
		ReadOnly: false,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if string(body) != `{"ok":true}` {
		t.Fatalf("body = %q", string(body))
	}
	if bodyReadErr != nil {
		t.Fatalf("request body read error = %v", bodyReadErr)
	}

	contractIssues := make([]string, 0, 4)
	if gotTenant != "tenant-a" {
		contractIssues = append(contractIssues, fmt.Sprintf("X-Tenant-Id = %q", gotTenant))
	}
	if gotAccessKey != "" {
		contractIssues = append(contractIssues, fmt.Sprintf("legacy X-Access-Key = %q", gotAccessKey))
	}
	if gotTimestamp != "" {
		contractIssues = append(contractIssues, fmt.Sprintf("legacy X-Timestamp = %q", gotTimestamp))
	}
	if decodeErr != nil {
		contractIssues = append(contractIssues, fmt.Sprintf("Signature header %q not JSON: %v", gotSignatureRaw, decodeErr))
	}
	if len(contractIssues) > 0 {
		t.Fatalf("SDK header contract violations: %s", strings.Join(contractIssues, "; "))
	}
	if len(gotSignature) != 3 {
		t.Fatalf("Signature = %#v, want 3 fields", gotSignature)
	}
	if gotSignature["AccessKey"] != "ak" {
		t.Fatalf("AccessKey = %q", gotSignature["AccessKey"])
	}
	if gotSignature["RequestSign"] == "" || gotSignature["RequestID"] == "" {
		t.Fatalf("Signature = %#v, want non-empty RequestSign and RequestID", gotSignature)
	}
	if string(gotBody) != `{"name":"daily"}` {
		t.Fatalf("request body = %q", string(gotBody))
	}
}

func TestExecute_OmitsTenantHeaderWhenTenantIDEmpty(t *testing.T) {
	var gotTenant string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTenant = r.Header.Get("X-Tenant-Id")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	exec := console.Executor{
		Signer: signer.Signer{},
		Client: transport.Client{Doer: srv.Client()},
	}

	_, err := exec.Execute(context.Background(), inputs.RemoteOptions{
		TenantID:      "",
		Endpoint:      srv.URL,
		AK:            "ak",
		SK:            "sk",
		TargetVersion: "9.0.9.0",
	}, console.RequestSpec{
		Method:   http.MethodGet,
		Path:     "/api/sla/v1/templates",
		ReadOnly: true,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotTenant != "" {
		t.Fatalf("X-Tenant-Id = %q, want empty", gotTenant)
	}
}
