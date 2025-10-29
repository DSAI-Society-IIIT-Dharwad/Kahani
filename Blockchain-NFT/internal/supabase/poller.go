package supabase

import (
	"context"
	"errors"
	"sync"
	"time"

	"storytelling-blockchain/internal/types"
)

// walletSyncClient defines the Supabase capabilities the poller relies on.
type walletSyncClient interface {
	FetchUsersSince(ctx context.Context, since time.Time) ([]User, error)
	UpsertWallet(ctx context.Context, wallet types.Wallet) error
}

// walletRegistry exposes the manager functions used to avoid duplicate wallets.
type walletRegistry interface {
	GetWalletBySupabaseID(userID string) (types.Wallet, bool)
}

// walletGenerator captures the generator capability required by the poller.
type walletGenerator interface {
	GenerateWalletForUser(supabaseUserID string) (types.Wallet, error)
}

// walletStorage captures the ability to persist wallets on-chain.
type walletStorage interface {
	StoreWalletOnChain(wallet types.Wallet) (types.Transaction, error)
}

// Poller orchestrates Supabase user polling and wallet provisioning.
type Poller struct {
	client    walletSyncClient
	generator walletGenerator
	storage   walletStorage
	manager   walletRegistry
	interval  time.Duration

	mu        sync.Mutex
	lastCheck time.Time
}

// NewPoller wires together the Supabase client, wallet subsystems, and poll interval.
func NewPoller(client walletSyncClient, generator walletGenerator, storage walletStorage, manager walletRegistry, interval time.Duration) (*Poller, error) {
	if client == nil || generator == nil || storage == nil || manager == nil {
		return nil, errors.New("supabase: poller dependencies must not be nil")
	}

	if interval <= 0 {
		interval = 30 * time.Second
	}

	return &Poller{
		client:    client,
		generator: generator,
		storage:   storage,
		manager:   manager,
		interval:  interval,
		lastCheck: time.Unix(0, 0),
	}, nil
}

// PollNewUsers pulls recent Supabase users and provisions wallets for newcomers.
func (p *Poller) PollNewUsers(ctx context.Context) (int, error) {
	p.mu.Lock()
	since := p.lastCheck
	p.mu.Unlock()

	users, err := p.client.FetchUsersSince(ctx, since)
	if err != nil {
		return 0, err
	}

	created := 0
	var newest time.Time

	for _, user := range users {
		if _, exists := p.manager.GetWalletBySupabaseID(user.ID); exists {
			if user.CreatedAt.After(newest) {
				newest = user.CreatedAt
			}
			continue
		}

		wallet, err := p.generator.GenerateWalletForUser(user.ID)
		if err != nil {
			return created, err
		}

		if err := p.client.UpsertWallet(ctx, wallet); err != nil {
			return created, err
		}

		if _, err := p.storage.StoreWalletOnChain(wallet); err != nil {
			return created, err
		}

		created++

		if user.CreatedAt.After(newest) {
			newest = user.CreatedAt
		}
	}

	if !newest.IsZero() {
		p.mu.Lock()
		if newest.After(p.lastCheck) {
			p.lastCheck = newest
		}
		p.mu.Unlock()
	}

	return created, nil
}

// Interval returns the configured polling interval.
func (p *Poller) Interval() time.Duration {
	return p.interval
}
