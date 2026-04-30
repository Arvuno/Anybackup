package api_gateway_token_to_x_user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// Config holds the Traefik plugin configuration.
type Config struct {
	HeaderName          string `json:"headerName,omitempty"`
	AuthHeaderName      string `json:"authHeaderName,omitempty"`
	UserinfoURL         string `json:"userinfoUrl,omitempty"`
	TimeoutMilliseconds int    `json:"timeoutMilliseconds,omitempty"`
	MaxUserinfoBytes    int64  `json:"maxUserinfoBytes,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderName:          "X-User",
		AuthHeaderName:      "Authorization",
		TimeoutMilliseconds: 5000,
		MaxUserinfoBytes:    65536,
	}
}

type tokenToXUser struct {
	next             http.Handler
	headerName       string
	authHeaderName   string
	userinfoURL      string
	maxUserinfoBytes int64
	httpClient       *http.Client
}

// New creates a Traefik middleware handler.
func New(_ context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	if config == nil {
		config = CreateConfig()
	}
	if config.HeaderName == "" {
		return nil, fmt.Errorf("headerName cannot be empty")
	}
	if config.AuthHeaderName == "" {
		return nil, fmt.Errorf("authHeaderName cannot be empty")
	}
	if config.UserinfoURL == "" {
		return nil, fmt.Errorf("userinfoUrl cannot be empty")
	}
	if config.TimeoutMilliseconds <= 0 {
		config.TimeoutMilliseconds = 5000
	}
	if config.MaxUserinfoBytes <= 0 {
		config.MaxUserinfoBytes = 65536
	}

	return &tokenToXUser{
		next:             next,
		headerName:       config.HeaderName,
		authHeaderName:   config.AuthHeaderName,
		userinfoURL:      config.UserinfoURL,
		maxUserinfoBytes: config.MaxUserinfoBytes,
		httpClient: &http.Client{
			Timeout: time.Duration(config.TimeoutMilliseconds) * time.Millisecond,
		},
	}, nil
}

func (m *tokenToXUser) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	token, ok := bearerToken(req.Header.Get(m.authHeaderName))
	if !ok {
		http.Error(rw, "missing bearer token", http.StatusUnauthorized)
		return
	}

	user, err := m.userInfo(req, token)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	xUser, err := json.Marshal(user.Claims)
	if err != nil {
		http.Error(rw, "invalid user context", http.StatusUnauthorized)
		return
	}

	req.Header.Set(m.headerName, string(xUser))
	req.Header.Del(m.authHeaderName)

	m.next.ServeHTTP(rw, req)
}

type userInfo struct {
	Claims map[string]any
}

func (u *userInfo) Sub() string {
	sub, _ := stringClaim(u.Claims, "sub")
	return sub
}

func bearerToken(header string) (string, bool) {
	fields := strings.Fields(header)
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") || fields[1] == "" {
		return "", false
	}
	return fields[1], true
}

func (m *tokenToXUser) userInfo(req *http.Request, token string) (*userInfo, error) {
	return m.userInfoFromEndpoint(req, token)
}

func (m *tokenToXUser) userInfoFromEndpoint(original *http.Request, token string) (*userInfo, error) {
	request, err := http.NewRequestWithContext(original.Context(), http.MethodGet, m.userinfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid userinfo URL")
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set(m.authHeaderName, "Bearer "+token)
	applyForwardedHeaders(request, original)

	response, err := m.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("userinfo request failed")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("userinfo rejected token")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("userinfo request failed")
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, m.maxUserinfoBytes+1))
	if err != nil {
		return nil, fmt.Errorf("userinfo response read failed")
	}
	if int64(len(body)) > m.maxUserinfoBytes {
		return nil, fmt.Errorf("userinfo response too large")
	}

	var claims map[string]any
	if err := json.Unmarshal(body, &claims); err != nil {
		return nil, fmt.Errorf("invalid userinfo response")
	}

	return userInfoFromClaims(claims)
}

func applyForwardedHeaders(request *http.Request, original *http.Request) {
	host := original.Host
	if host == "" {
		host = original.Header.Get("Host")
	}
	if host != "" {
		request.Host = host
		request.Header.Set("X-Forwarded-Host", firstNonEmpty(original.Header.Get("X-Forwarded-Host"), host))
	}

	proto := firstNonEmpty(original.Header.Get("X-Forwarded-Proto"), requestScheme(original))
	request.Header.Set("X-Forwarded-Proto", proto)
	request.Header.Set("X-Forwarded-Port", firstNonEmpty(original.Header.Get("X-Forwarded-Port"), requestPort(host, proto)))

	if forwardedFor := original.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		request.Header.Set("X-Forwarded-For", forwardedFor)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func requestScheme(req *http.Request) string {
	if req.TLS != nil {
		return "https"
	}
	return "http"
}

func requestPort(host string, proto string) string {
	if _, port, err := net.SplitHostPort(host); err == nil {
		return port
	}
	if proto == "https" {
		return "443"
	}
	return "80"
}

func decodeClaims(token string) (map[string]any, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("token must contain header and payload")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		payload, err = base64.URLEncoding.DecodeString(parts[1])
		if err != nil {
			return nil, err
		}
	}

	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	return claims, nil
}

func userInfoFromClaims(claims map[string]any) (*userInfo, error) {
	sub, ok := stringClaim(claims, "sub")
	if !ok {
		return nil, fmt.Errorf("sub claim is required")
	}

	user := &userInfo{Claims: claims}
	if roles := rolesFromClaims(claims); len(roles) > 0 {
		user.Claims["roles"] = roles
	}
	user.Claims["sub"] = sub

	return user, nil
}

func stringClaim(claims map[string]any, name string) (string, bool) {
	value, ok := claims[name].(string)
	if !ok || value == "" {
		return "", false
	}
	return value, true
}

func rolesFromClaims(claims map[string]any) []string {
	seen := map[string]bool{}
	var roles []string

	addRoles := func(values any) {
		items, ok := values.([]any)
		if !ok {
			return
		}
		for _, item := range items {
			role, ok := item.(string)
			if !ok || role == "" || seen[role] {
				continue
			}
			seen[role] = true
			roles = append(roles, role)
		}
	}

	if realmAccess, ok := claims["realm_access"].(map[string]any); ok {
		addRoles(realmAccess["roles"])
	}

	if resourceAccess, ok := claims["resource_access"].(map[string]any); ok {
		for _, access := range resourceAccess {
			accessMap, ok := access.(map[string]any)
			if !ok {
				continue
			}
			addRoles(accessMap["roles"])
		}
	}

	return roles
}
