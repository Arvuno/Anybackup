package api_gateway_token_to_x_user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTokenToXUserSetsUserinfoHeaderFromBearerToken(t *testing.T) {
	userinfoServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatalf("unexpected authorization header: %s", req.Header.Get("Authorization"))
		}
		if req.Host != "example.test" {
			t.Fatalf("unexpected userinfo host: %s", req.Host)
		}
		if req.Header.Get("X-Forwarded-Host") != "example.test" {
			t.Fatalf("unexpected X-Forwarded-Host: %s", req.Header.Get("X-Forwarded-Host"))
		}
		if req.Header.Get("X-Forwarded-Proto") != "http" {
			t.Fatalf("unexpected X-Forwarded-Proto: %s", req.Header.Get("X-Forwarded-Proto"))
		}
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{
			"sub": "user-001",
			"preferred_username": "admin",
			"email": "admin@example.com",
			"email_verified": true,
			"roles": ["backup_admin", "offline_access"]
		}`))
	}))
	defer userinfoServer.Close()

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		xUser := req.Header.Get("X-User")
		if xUser == "" {
			t.Fatal("expected X-User header")
		}
		if req.Header.Get("Authorization") != "" {
			t.Fatalf("expected Authorization header to be removed, got %q", req.Header.Get("Authorization"))
		}
		if req.Header.Get("X-Request-User") != "" {
			t.Fatalf("expected X-Request-User to remain unset, got %q", req.Header.Get("X-Request-User"))
		}
		if req.Header.Get("X-Auth-Request-User") != "" {
			t.Fatalf("expected X-Auth-Request-User to remain unset, got %q", req.Header.Get("X-Auth-Request-User"))
		}

		var user map[string]any
		if err := json.Unmarshal([]byte(xUser), &user); err != nil {
			t.Fatalf("expected X-User to contain JSON: %v", err)
		}

		if user["sub"] != "user-001" {
			t.Fatalf("unexpected sub: %v", user["sub"])
		}
		if user["preferred_username"] != "admin" {
			t.Fatalf("unexpected preferred_username: %v", user["preferred_username"])
		}
		if user["email"] != "admin@example.com" {
			t.Fatalf("unexpected email: %v", user["email"])
		}

		roles, ok := user["roles"].([]any)
		if !ok || len(roles) != 2 || roles[0] != "backup_admin" || roles[1] != "offline_access" {
			t.Fatalf("unexpected roles: %#v", user["roles"])
		}

		rw.WriteHeader(http.StatusNoContent)
	})

	config := CreateConfig()
	config.UserinfoURL = userinfoServer.URL
	handler, err := New(context.Background(), next, config, "token-to-x-user")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "http://example.test/resource", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", response.Code)
	}
}

func TestTokenToXUserRejectsUserinfoWithoutSubject(t *testing.T) {
	userinfoServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(`{"preferred_username":"missing-sub"}`))
	}))
	defer userinfoServer.Close()

	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler should not be called")
	})

	config := CreateConfig()
	config.UserinfoURL = userinfoServer.URL
	handler, err := New(context.Background(), next, config, "token-to-x-user")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "http://example.test/resource", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", response.Code)
	}
}

func TestTokenToXUserRejectsMissingBearerToken(t *testing.T) {
	userinfoServer := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("userinfo endpoint should not be called")
	}))
	defer userinfoServer.Close()

	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next handler should not be called")
	})

	config := CreateConfig()
	config.UserinfoURL = userinfoServer.URL
	handler, err := New(context.Background(), next, config, "token-to-x-user")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "http://example.test/resource", nil))

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", response.Code)
	}
}

func jwtWithPayload(payload map[string]any) string {
	header := encodeJWTPart(map[string]any{"alg": "none", "typ": "JWT"})
	body := encodeJWTPart(payload)
	return header + "." + body + "."
}

func encodeJWTPart(value map[string]any) string {
	raw, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(raw)
}
