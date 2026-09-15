package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	users := NewUserStore()
	if err := users.Register("alice", "s3cr3t-password"); err != nil {
		t.Fatalf("failed to register user: %v", err)
	}
	tokens, err := NewTokenIssuer(strings.Repeat("x", 32), time.Minute)
	if err != nil {
		t.Fatalf("failed to create token issuer: %v", err)
	}
	return NewServer(users, tokens)
}

func TestHealthz(t *testing.T) {
	server := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestLoginSuccessAndValidate(t *testing.T) {
	server := newTestServer(t)

	body, _ := json.Marshal(loginRequest{Username: "alice", Password: "s3cr3t-password"})
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var loginResp loginResponse
	if err := json.NewDecoder(rec.Body).Decode(&loginResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if loginResp.Token == "" {
		t.Fatal("expected non-empty token")
	}

	validateReq := httptest.NewRequest(http.MethodGet, "/validate", nil)
	validateReq.Header.Set("Authorization", "Bearer "+loginResp.Token)
	validateRec := httptest.NewRecorder()
	server.ServeHTTP(validateRec, validateReq)

	if validateRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", validateRec.Code)
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	server := newTestServer(t)

	body, _ := json.Marshal(loginRequest{Username: "alice", Password: "wrong-password"})
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(body)))
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestValidateMissingToken(t *testing.T) {
	server := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/validate", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}
