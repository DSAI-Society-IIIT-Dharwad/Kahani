package supabase

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type verifierStub struct {
	userID string
	err    error
}

func (v verifierStub) VerifyToken(ctx context.Context, token string) (string, error) {
	return v.userID, v.err
}

func TestAuthMiddlewareSuccess(t *testing.T) {
	middleware, err := NewAuthMiddleware(verifierStub{userID: "user-123"})
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	didCall := false
	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		didCall = true
		userID, ok := UserIDFromContext(r.Context())
		if !ok || userID != "user-123" {
			t.Fatalf("expected user id to be injected, got %q", userID)
		}
	}))

	handler.ServeHTTP(w, req)

	if !didCall {
		t.Fatalf("expected next handler to execute")
	}

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", w.Code)
	}
}

func TestAuthMiddlewareMissingHeader(t *testing.T) {
	middleware, _ := NewAuthMiddleware(verifierStub{userID: "user"})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	middleware.Wrap(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing header, got %d", w.Code)
	}
}

func TestAuthMiddlewareVerificationFailure(t *testing.T) {
	middleware, _ := NewAuthMiddleware(verifierStub{userID: ""})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	middleware.Wrap(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when verification fails, got %d", w.Code)
	}
}

func TestUserIDFromContextAbsent(t *testing.T) {
	if _, ok := UserIDFromContext(context.Background()); ok {
		t.Fatalf("expected missing user id to return false")
	}
}
