package supabase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sort"
	"storytelling-blockchain/internal/types"
	"time"
)

// Config encapsulates the credentials required for Supabase access.
type Config struct {
	URL            string
	AnonKey        string
	ServiceRoleKey string
	HTTPClient     *http.Client
}

// Client talks to Supabase Auth and Admin endpoints.
type Client struct {
	baseURL        *url.URL
	anonKey        string
	serviceRoleKey string
	httpClient     *http.Client
}

// User represents the subset of Supabase auth user fields used by the poller.
type User struct {
	ID        string
	CreatedAt time.Time
}

// NewClient validates configuration and returns a Supabase client instance.
func NewClient(cfg Config) (*Client, error) {
	if cfg.URL == "" {
		return nil, errors.New("supabase: missing url")
	}

	parsed, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("supabase: invalid url: %w", err)
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &Client{
		baseURL:        parsed,
		anonKey:        cfg.AnonKey,
		serviceRoleKey: cfg.ServiceRoleKey,
		httpClient:     httpClient,
	}, nil
}

func (c *Client) endpoint(parts ...string) string {
	joined := path.Join(parts...)
	return c.baseURL.ResolveReference(&url.URL{Path: joined}).String()
}

// VerifyToken checks the provided JWT via Supabase's auth user endpoint and returns the user ID.
func (c *Client) VerifyToken(ctx context.Context, jwtToken string) (string, error) {
	if jwtToken == "" {
		return "", errors.New("supabase: empty jwt token")
	}

	if c.anonKey == "" {
		return "", errors.New("supabase: anon key required for verification")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint("auth", "v1", "user"), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("apikey", c.anonKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("supabase: verify token failed with status %d", resp.StatusCode)
	}

	var payload struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("supabase: decode verify response failed: %w", err)
	}

	if payload.ID == "" {
		return "", errors.New("supabase: missing user id in response")
	}

	return payload.ID, nil
}

// FetchUsersSince retrieves auth users created after the provided timestamp using the admin endpoint.
func (c *Client) FetchUsersSince(ctx context.Context, since time.Time) ([]User, error) {
	if c.serviceRoleKey == "" {
		return nil, errors.New("supabase: service role key required for admin operations")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint("auth", "v1", "admin", "users"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.serviceRoleKey)
	req.Header.Set("apikey", c.serviceRoleKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase: fetch users failed with status %d", resp.StatusCode)
	}

	var payload struct {
		Users []struct {
			ID        string `json:"id"`
			CreatedAt string `json:"created_at"`
		} `json:"users"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("supabase: decode admin users failed: %w", err)
	}

	var result []User
	for _, user := range payload.Users {
		createdAt, err := time.Parse(time.RFC3339, user.CreatedAt)
		if err != nil {
			continue
		}
		if createdAt.After(since) {
			result = append(result, User{ID: user.ID, CreatedAt: createdAt})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result, nil
}

// UpsertWallet synchronizes a wallet record into Supabase using the REST API.
func (c *Client) UpsertWallet(ctx context.Context, wallet types.Wallet) error {
	if c.serviceRoleKey == "" {
		return errors.New("supabase: service role key required for admin operations")
	}

	if wallet.SupabaseUserID == "" {
		return errors.New("supabase: wallet missing supabase user id")
	}

	endpoint := c.baseURL.ResolveReference(&url.URL{
		Path:     path.Join("rest", "v1", "wallets"),
		RawQuery: "on_conflict=supabase_user_id",
	})

	payload := []map[string]interface{}{
		{
			"supabase_user_id":      wallet.SupabaseUserID,
			"address":               wallet.Address,
			"public_key":            wallet.PublicKey,
			"private_key_encrypted": wallet.PrivateKeyEncrypted,
			"created_at":            wallet.CreatedAt,
			"block_index":           wallet.BlockIndex,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("supabase: encode wallet payload failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "resolution=merge-duplicates")
	req.Header.Set("Authorization", "Bearer "+c.serviceRoleKey)
	req.Header.Set("apikey", c.serviceRoleKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("supabase: upsert wallet failed with status %d", resp.StatusCode)
	}

	return nil
}

// BuildSupabaseInstructions returns high-level setup steps for configuring Supabase.
func BuildSupabaseInstructions() []string {
	return []string{
		"Create a new Supabase project and note the project URL",
		"Generate anon and service role keys from Project Settings → API",
		"Ensure the service role key remains server-side only",
		"Enable Email or OAuth providers as needed for your frontend",
		"(Optional) Restrict access to the admin users endpoint via RLS policies",
	}
}
