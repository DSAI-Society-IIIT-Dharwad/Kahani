package supabase

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"storytelling-blockchain/internal/types"
)

func TestClientVerifyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/v1/user":
			if r.Header.Get("Authorization") != "Bearer good-token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"id": "user-1"})
		case "/auth/v1/admin/users":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"users": []map[string]string{{
					"id":         "user-1",
					"created_at": time.Now().Add(-time.Hour).Format(time.RFC3339),
				}},
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		URL:            ts.URL,
		AnonKey:        "anon",
		ServiceRoleKey: "service",
		HTTPClient:     ts.Client(),
	})
	if err != nil {
		t.Fatalf("client init failed: %v", err)
	}

	ctx := context.Background()
	userID, err := client.VerifyToken(ctx, "good-token")
	if err != nil {
		t.Fatalf("verify token failed: %v", err)
	}

	if userID != "user-1" {
		t.Fatalf("unexpected user id: %s", userID)
	}

	if _, err := client.VerifyToken(ctx, "bad-token"); err == nil {
		t.Fatalf("expected error for bad token")
	}
}

func TestFetchUsersSinceFilters(t *testing.T) {
	now := time.Now()
	earlier := now.Add(-2 * time.Hour)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"users": []map[string]string{
				{"id": "user-old", "created_at": earlier.Format(time.RFC3339)},
				{"id": "user-new", "created_at": now.Format(time.RFC3339)},
			},
		})
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		URL:            ts.URL,
		ServiceRoleKey: "service",
		HTTPClient:     ts.Client(),
	})
	if err != nil {
		t.Fatalf("client init failed: %v", err)
	}

	ctx := context.Background()
	users, err := client.FetchUsersSince(ctx, earlier.Add(time.Minute))
	if err != nil {
		t.Fatalf("fetch users failed: %v", err)
	}

	if len(users) != 1 || users[0].ID != "user-new" {
		t.Fatalf("expected only new user to be returned, got %+v", users)
	}

	if len(BuildSupabaseInstructions()) == 0 {
		t.Fatalf("expected supabase instructions to be populated")
	}
}

func TestUpsertWallet(t *testing.T) {
	var (
		receivedPath    string
		receivedQuery   string
		receivedHeaders http.Header
		receivedBody    []byte
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		receivedQuery = r.URL.RawQuery
		receivedHeaders = r.Header.Clone()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		receivedBody = body
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	client, err := NewClient(Config{
		URL:            ts.URL,
		ServiceRoleKey: "service",
		HTTPClient:     ts.Client(),
	})
	if err != nil {
		t.Fatalf("client init failed: %v", err)
	}

	wallet := types.Wallet{
		SupabaseUserID:      "f95c1d72-7896-4b2a-a9e4-1d8facf8a0d1",
		Address:             "0x123",
		PublicKey:           "pub",
		PrivateKeyEncrypted: "enc",
		CreatedAt:           42,
		BlockIndex:          -1,
	}

	if err := client.UpsertWallet(context.Background(), wallet); err != nil {
		t.Fatalf("upsert wallet failed: %v", err)
	}

	if receivedPath != "/rest/v1/wallets" {
		t.Fatalf("unexpected path: %s", receivedPath)
	}

	if receivedQuery != "on_conflict=supabase_user_id" {
		t.Fatalf("unexpected query: %s", receivedQuery)
	}

	if receivedHeaders.Get("Prefer") != "resolution=merge-duplicates" {
		t.Fatalf("prefer header missing")
	}

	var payload []map[string]interface{}
	if err := json.Unmarshal(receivedBody, &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	if len(payload) != 1 {
		t.Fatalf("expected single payload object, got %d", len(payload))
	}

	if payload[0]["supabase_user_id"] != wallet.SupabaseUserID {
		t.Fatalf("unexpected wallet payload: %+v", payload[0])
	}
}
