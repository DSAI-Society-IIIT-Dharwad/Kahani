package supabase

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// TokenVerifier is implemented by clients capable of verifying Supabase JWT tokens.
type TokenVerifier interface {
	VerifyToken(ctx context.Context, jwtToken string) (string, error)
}

// AuthMiddleware handles Supabase JWT verification and context propagation.
type AuthMiddleware struct {
	verifier TokenVerifier
}

type contextKey string

const userIDContextKey contextKey = "supabase_user_id"

// ErrMissingAuthorization indicates the request lacked a valid Authorization header.
var ErrMissingAuthorization = errors.New("supabase: missing or invalid authorization header")

// NewAuthMiddleware constructs an AuthMiddleware from the provided verifier.
func NewAuthMiddleware(verifier TokenVerifier) (*AuthMiddleware, error) {
	if verifier == nil {
		return nil, errors.New("supabase: token verifier is nil")
	}
	return &AuthMiddleware{verifier: verifier}, nil
}

// Wrap adapts the middleware to standard http handlers.
func (a *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	if a == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractBearerToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, ErrMissingAuthorization.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := a.verifier.VerifyToken(r.Context(), token)
		if err != nil || userID == "" {
			http.Error(w, "supabase: token verification failed", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", ErrMissingAuthorization
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrMissingAuthorization
	}

	if strings.TrimSpace(parts[1]) == "" {
		return "", ErrMissingAuthorization
	}

	return parts[1], nil
}

// UserIDFromContext extracts the Supabase user ID from the request context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(userIDContextKey).(string)
	return val, ok && val != ""
}
