package signer_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/anybackup-ai/Anybackup/CLI/internal/platform/signer"
)

var lowerHexPattern = regexp.MustCompile(`^[0-9a-f]+$`)

func decodeSignatureMap(t *testing.T, raw string) map[string]string {
	t.Helper()
	if raw == "" {
		t.Fatal("Signature header is empty")
	}
	got := map[string]string{}
	if err := json.Unmarshal([]byte(raw), &got); err != nil {
		t.Fatalf("Signature header is not valid JSON: %v", err)
	}
	return got
}

// expectedRequestSign mirrors the SDK-style normalization (method+host+path+query/body) hashed HMAC-SHA256.
func expectedRequestSign(t *testing.T, req *http.Request, body []byte, sk string) string {
	t.Helper()

	var payload bytes.Buffer
	payload.WriteString(strings.ToUpper(req.Method))
	payload.WriteString(req.URL.Host)
	payload.WriteString(req.URL.EscapedPath())

	if req.URL.RawQuery != "" || len(body) > 0 {
		payload.WriteByte('?')
		if req.URL.RawQuery != "" {
			query, err := url.PathUnescape(req.URL.RawQuery)
			if err != nil {
				t.Fatalf("PathUnescape(%q) error = %v", req.URL.RawQuery, err)
			}
			payload.WriteString(query)
		}
	}
	payload.Write(body)

	mac := hmac.New(sha256.New, []byte(sk))
	_, _ = mac.Write(payload.Bytes())
	return hex.EncodeToString(mac.Sum(nil))
}

func TestSigner_RequestSignMatchesSDKStyleContract(t *testing.T) {
	cases := []struct {
		name   string
		method string
		rawURL string
		body   []byte
	}{
		{
			name:   "path only",
			method: http.MethodGet,
			rawURL: "https://127.0.0.1:9600/api/sla/v1/templates",
		},
		{
			name:   "path and query",
			method: http.MethodGet,
			rawURL: "https://127.0.0.1:9600/api/sla/v1/templates?scope=a%2Fb",
		},
		{
			name:   "path and body",
			method: http.MethodPost,
			rawURL: "https://127.0.0.1:9600/backupmgm/v1/file/recovery",
			body:   []byte(`{"name":"daily"}`),
		},
		{
			name:   "path query and body",
			method: http.MethodPost,
			rawURL: "https://127.0.0.1:9600/api/file%20recovery?scope=a%2Fb",
			body:   []byte(`{"name":"daily"}`),
		},
		{
			name:   "escaped path",
			method: http.MethodGet,
			rawURL: "https://127.0.0.1:9600/api/%E4%B8%AD%E6%96%87/file%20recovery",
		},
		{
			name:   "query canonicalization mix",
			method: http.MethodGet,
			rawURL: "https://127.0.0.1:9600/api/file%20recovery?scope=a%2Fb&filter=foo+bar&filter=baz%20qux&special=1%2B2",
		},
		{
			name:   "query encoded delimiters",
			method: http.MethodGet,
			rawURL: "https://127.0.0.1:9600/api/file%20recovery?name=a%26b%3Dc&desc=hello%26world%3D",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.rawURL, bytes.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}

			err = signer.Signer{}.Sign(req, signer.Credentials{
				TenantID: "tenant-a",
				AK:       "ak",
				SK:       "sk",
			})
			if err != nil {
				t.Fatalf("Sign() error = %v", err)
			}

			if req.Header.Get("X-Tenant-Id") != "tenant-a" {
				t.Fatalf("X-Tenant-Id = %q", req.Header.Get("X-Tenant-Id"))
			}
			if got := req.Header.Get("X-Access-Key"); got != "" {
				t.Fatalf("X-Access-Key = %q, want empty", got)
			}
			got := decodeSignatureMap(t, req.Header.Get("Signature"))
			if len(got) != 3 {
				t.Fatalf("Signature fields = %#v, want exactly 3 keys", got)
			}
			if got["AccessKey"] != "ak" {
				t.Fatalf("AccessKey = %q", got["AccessKey"])
			}

			wantSign := expectedRequestSign(t, req, tc.body, "sk")
			if got["RequestSign"] != wantSign {
				t.Fatalf("RequestSign = %q, want %q", got["RequestSign"], wantSign)
			}
			if !lowerHexPattern.MatchString(got["RequestID"]) {
				t.Fatalf("RequestID = %q, want lowercase hex", got["RequestID"])
			}
		})
	}
}

func TestSigner_PreservesRequestBodyAfterSigning(t *testing.T) {
	t.Run("non-empty body is preserved", func(t *testing.T) {
		body := []byte(`{"objectId":"obj-1","name":"daily"}`)
		req, err := http.NewRequest(
			http.MethodPost,
			"https://127.0.0.1:9600/backupmgm/v1/file/recovery?scope=a%2Fb",
			bytes.NewReader(body),
		)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		err = signer.Signer{}.Sign(req, signer.Credentials{
			TenantID: "tenant-a",
			AK:       "ak",
			SK:       "sk",
		})
		if err != nil {
			t.Fatalf("Sign() error = %v", err)
		}

		got, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("ReadAll(req.Body) error = %v", err)
		}
		if !bytes.Equal(got, body) {
			t.Fatalf("body = %q, want %q", string(got), string(body))
		}
		if got := req.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}
	})

	t.Run("empty body preserves content-type", func(t *testing.T) {
		req, err := http.NewRequest(
			http.MethodPost,
			"https://127.0.0.1:9600/backupmgm/v1/file/recovery",
			bytes.NewReader([]byte{}),
		)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		err = signer.Signer{}.Sign(req, signer.Credentials{
			TenantID: "tenant-a",
			AK:       "ak",
			SK:       "sk",
		})
		if err != nil {
			t.Fatalf("Sign() error = %v", err)
		}

		if got := req.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}

		gotBody, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("ReadAll(req.Body) error = %v", err)
		}
		if len(gotBody) != 0 {
			t.Fatalf("body = %q, want empty", string(gotBody))
		}
	})
}

func TestSigner_UsesSignPathOverrideAndRemovesInternalHeader(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://127.0.0.1:9600/storageresmgm/v1/storage-svc-id/pool?index=0&count=30",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Foundation-Cli-Sign-Path", "/storageresmgm/v1/pool")

	err = signer.Signer{}.Sign(req, signer.Credentials{
		TenantID: "tenant-a",
		AK:       "ak",
		SK:       "sk",
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	if got := req.Header.Get("X-Foundation-Cli-Sign-Path"); got != "" {
		t.Fatalf("X-Foundation-Cli-Sign-Path = %q, want empty", got)
	}

	got := decodeSignatureMap(t, req.Header.Get("Signature"))
	overrideURL := *req.URL
	overrideURL.Path = "/storageresmgm/v1/pool"
	overrideReq := *req
	overrideReq.URL = &overrideURL
	wantSign := expectedRequestSign(t, &overrideReq, nil, "sk")
	if got["RequestSign"] != wantSign {
		t.Fatalf("RequestSign = %q, want %q", got["RequestSign"], wantSign)
	}
}

func TestSigner_OmitsTenantHeaderWhenTenantIDEmpty(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://127.0.0.1:9600/api/sla/v1/templates",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Tenant-Id", "legacy-tenant")

	err = signer.Signer{}.Sign(req, signer.Credentials{
		TenantID: "",
		AK:       "ak",
		SK:       "sk",
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	if got := req.Header.Get("X-Tenant-Id"); got != "" {
		t.Fatalf("X-Tenant-Id = %q, want empty", got)
	}
}
