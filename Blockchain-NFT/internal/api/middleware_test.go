package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"storytelling-blockchain/internal/supabase"
)

func TestMiddlewareWrapAppliesCORS(t *testing.T) {
	mw := NewMiddleware(MiddlewareConfig{
		AllowedOrigins: []string{"http://example.com"},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()

	mw.Wrap(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://example.com" {
		t.Fatalf("expected CORS header to be set, got %q", got)
	}
}

func TestMiddlewareWrapAuthenticatedWithRateLimit(t *testing.T) {
	stub := &verifierStub{token: "good"}
	auth, err := supabase.NewAuthMiddleware(stub)
	if err != nil {
		t.Fatalf("failed to create auth middleware: %v", err)
	}

	mw := NewMiddleware(MiddlewareConfig{
		Auth:      auth,
		RateLimit: rate.Every(time.Hour),
		Burst:     1,
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer good")
	w := httptest.NewRecorder()

	called := false
	handler := mw.WrapAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := supabase.UserIDFromContext(r.Context()); !ok {
			t.Fatalf("expected user id in context")
		}
		called = true
	}))

	handler.ServeHTTP(w, req)

	if !called {
		t.Fatalf("expected downstream handler to be invoked")
	}

	if !stub.used {
		t.Fatalf("expected token verifier to be called")
	}

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from downstream handler, got %d", w.Code)
	}

	// Trigger rate limit by setting limiter to zero allowance.
	mw.limiter = rate.NewLimiter(0, 0)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req)
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 when rate limit exceeded, got %d", w2.Code)
	}
}

type verifierStub struct {
	token string
	used  bool
}

func (v *verifierStub) VerifyToken(_ context.Context, token string) (string, error) {
	if token != v.token {
		return "", errors.New("invalid token")
	}
	v.used = true
	return "user-123", nil
}
