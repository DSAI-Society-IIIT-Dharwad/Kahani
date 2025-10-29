package supabase

import (
	"context"
	"testing"
	"time"

	"storytelling-blockchain/internal/blockchain"
	"storytelling-blockchain/internal/types"
	"storytelling-blockchain/internal/wallet"
)

type fakeSupabaseClient struct {
	users   []User
	upserts []types.Wallet
}

func (f *fakeSupabaseClient) FetchUsersSince(ctx context.Context, since time.Time) ([]User, error) {
	return f.users, nil
}

func (f *fakeSupabaseClient) UpsertWallet(ctx context.Context, wallet types.Wallet) error {
	f.upserts = append(f.upserts, wallet)
	return nil
}

func TestPollerCreatesWallets(t *testing.T) {
	originalNow := types.NowUnix
	types.NowUnix = func() int64 { return 777 }
	t.Cleanup(func() { types.NowUnix = originalNow })

	chain := blockchain.NewBlockchain()
	generator, err := wallet.NewGenerator("passphrase")
	if err != nil {
		t.Fatalf("generator init failed: %v", err)
	}

	storage, err := wallet.NewStorage(chain)
	if err != nil {
		t.Fatalf("storage init failed: %v", err)
	}

	manager, err := wallet.NewManager(chain, "passphrase")
	if err != nil {
		t.Fatalf("manager init failed: %v", err)
	}

	client := &fakeSupabaseClient{
		users: []User{
			{ID: "user-a", CreatedAt: time.Now().Add(-time.Minute)},
			{ID: "user-b", CreatedAt: time.Now()},
		},
	}

	poller, err := NewPoller(client, generator, storage, manager, 30*time.Second)
	if err != nil {
		t.Fatalf("poller init failed: %v", err)
	}

	count, err := poller.PollNewUsers(context.Background())
	if err != nil {
		t.Fatalf("polling failed: %v", err)
	}

	if count != 2 {
		t.Fatalf("expected two wallets created, got %d", count)
	}

	state := chain.State()
	if len(state.WalletRegistry) != 2 {
		t.Fatalf("expected two wallets in state, got %d", len(state.WalletRegistry))
	}

	if len(client.upserts) != 2 {
		t.Fatalf("expected two supabase wallet upserts, got %d", len(client.upserts))
	}

	// Second poll should be idempotent once lastCheck is updated.
	count, err = poller.PollNewUsers(context.Background())
	if err != nil {
		t.Fatalf("polling failed: %v", err)
	}

	if count != 0 {
		t.Fatalf("expected no new wallets on second poll, got %d", count)
	}
}
