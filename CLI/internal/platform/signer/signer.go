package signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Credentials struct {
	TenantID string
	AK       string
	SK       string
}

type Signer struct{}

const signPathOverrideHeader = "X-Foundation-Cli-Sign-Path"

type signatureHeader struct {
	AccessKey   string `json:"AccessKey"`
	RequestSign string `json:"RequestSign"`
	RequestID   string `json:"RequestID"`
}

func (Signer) Sign(req *http.Request, creds Credentials) error {
	body, err := readAndRestoreBody(req)
	if err != nil {
		return err
	}

	stringToSign, err := buildStringToSign(req, body)
	if err != nil {
		return err
	}

	signature, err := json.Marshal(signatureHeader{
		AccessKey:   creds.AK,
		RequestSign: makeHMAC(stringToSign, []byte(creds.SK), sha256.New),
		RequestID:   makeHMAC([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)), []byte(creds.SK), sha1.New),
	})
	if err != nil {
		return err
	}

	if tenantID := strings.TrimSpace(creds.TenantID); tenantID == "" {
		req.Header.Del("X-Tenant-Id")
	} else {
		req.Header.Set("X-Tenant-Id", tenantID)
	}
	req.Header.Del("X-Access-Key")
	req.Header.Del("X-Timestamp")
	req.Header.Del(signPathOverrideHeader)
	req.Header.Set("Signature", string(signature))
	debugSign(req, body, stringToSign)
	return nil
}

func buildStringToSign(req *http.Request, body []byte) ([]byte, error) {
	var payload bytes.Buffer
	payload.WriteString(strings.ToUpper(req.Method))
	payload.WriteString(signatureHost(req))
	payload.WriteString(signaturePath(req))

	if req.URL.RawQuery != "" || len(body) > 0 {
		payload.WriteByte('?')
		if req.URL.RawQuery != "" {
			query, err := url.PathUnescape(req.URL.RawQuery)
			if err != nil {
				return nil, err
			}
			payload.WriteString(query)
		}
		payload.Write(body)
	}

	return payload.Bytes(), nil
}

func signatureHost(req *http.Request) string {
	if os.Getenv("FOUNDATION_CLI_SIGN_HOST_MODE") == "hostname_only" {
		return req.URL.Hostname()
	}
	return req.URL.Host
}

func signaturePath(req *http.Request) string {
	if override := req.Header.Get(signPathOverrideHeader); override != "" {
		return override
	}
	return req.URL.EscapedPath()
}

func readAndRestoreBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	if err := req.Body.Close(); err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}
	return body, nil
}

func makeHMAC(payload []byte, key []byte, newHash func() hash.Hash) string {
	mac := hmac.New(newHash, key)
	_, _ = mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func debugSign(req *http.Request, body []byte, stringToSign []byte) {
	if os.Getenv("FOUNDATION_CLI_DEBUG_SIGN") == "" {
		return
	}

	_, _ = fmt.Fprintf(os.Stderr, "[sign-debug] method=%s\n", req.Method)
	_, _ = fmt.Fprintf(os.Stderr, "[sign-debug] url=%s\n", req.URL.String())
	_, _ = fmt.Fprintf(os.Stderr, "[sign-debug] host=%s\n", req.URL.Host)
	_, _ = fmt.Fprintf(os.Stderr, "[sign-debug] escaped_path=%s\n", req.URL.EscapedPath())
	_, _ = fmt.Fprintf(os.Stderr, "[sign-debug] raw_query=%s\n", req.URL.RawQuery)
	_, _ = fmt.Fprintf(os.Stderr, "[sign-debug] body=%s\n", string(body))
	_, _ = fmt.Fprintf(os.Stderr, "[sign-debug] string_to_sign=%s\n", string(stringToSign))
	_, _ = fmt.Fprintf(os.Stderr, "[sign-debug] x_tenant_id=%s\n", req.Header.Get("X-Tenant-Id"))
	_, _ = fmt.Fprintf(os.Stderr, "[sign-debug] signature=%s\n", req.Header.Get("Signature"))
}
