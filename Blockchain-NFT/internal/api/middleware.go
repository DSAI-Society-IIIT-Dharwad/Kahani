package api

import (
	"net/http"
	"strings"

	"golang.org/x/time/rate"

	"storytelling-blockchain/internal/supabase"
)

// MiddlewareConfig captures configuration for API middleware features.
type MiddlewareConfig struct {
	Auth           *supabase.AuthMiddleware
	RateLimit      rate.Limit
	Burst          int
	AllowedOrigins []string
}

// Middleware applies CORS headers, rate limiting, and optional Supabase auth.
type Middleware struct {
	auth           *supabase.AuthMiddleware
	limiter        *rate.Limiter
	allowedOrigins []string
}

// NewMiddleware constructs a Middleware instance using the provided configuration.
func NewMiddleware(cfg MiddlewareConfig) *Middleware {
	var limiter *rate.Limiter
	if cfg.RateLimit > 0 {
		burst := cfg.Burst
		if burst <= 0 {
			burst = 1
		}
		limiter = rate.NewLimiter(cfg.RateLimit, burst)
	}

	return &Middleware{
		auth:           cfg.Auth,
		limiter:        limiter,
		allowedOrigins: cfg.AllowedOrigins,
	}
}

// Wrap applies CORS and rate limiting (but not auth) to the handler chain.
func (m *Middleware) Wrap(next http.Handler) http.Handler {
	if m == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.applyCORS(w, r)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if m.limiter != nil && !m.limiter.Allow() {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// WrapAuthenticated applies auth, CORS, and rate limiting to the chain.
func (m *Middleware) WrapAuthenticated(next http.Handler) http.Handler {
	if m == nil {
		return next
	}

	wrapped := m.Wrap(next)
	if m.auth != nil {
		wrapped = m.auth.Wrap(wrapped)
	}
	return wrapped
}

// AuthMiddleware exposes the underlying Supabase auth middleware.
func (m *Middleware) AuthMiddleware() *supabase.AuthMiddleware {
	if m == nil {
		return nil
	}
	return m.auth
}

func (m *Middleware) applyCORS(w http.ResponseWriter, r *http.Request) {
	if len(m.allowedOrigins) == 0 {
		return
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = m.allowedOrigins[0]
	}

	if !m.originAllowed(origin) {
		return
	}

	headers := w.Header()
	headers.Set("Access-Control-Allow-Origin", origin)
	headers.Set("Vary", "Origin")
	headers.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	headers.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

func (m *Middleware) originAllowed(origin string) bool {
	if len(m.allowedOrigins) == 0 {
		return false
	}

	for _, allowed := range m.allowedOrigins {
		if strings.EqualFold(allowed, origin) {
			return true
		}
	}

	return false
}

// OriginAllowed reports whether the provided origin is permitted.
func (m *Middleware) OriginAllowed(origin string) bool {
	if m == nil {
		return true
	}

	if len(m.allowedOrigins) == 0 {
		return true
	}

	if origin == "" {
		return false
	}

	return m.originAllowed(origin)
}
