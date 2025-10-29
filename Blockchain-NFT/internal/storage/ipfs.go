package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"sync"

	shell "github.com/ipfs/go-ipfs-api"

	"storytelling-blockchain/pkg/utils"
)

// IPFSClient defines the methods required by higher level components.
type IPFSClient interface {
	UploadBytes(data []byte) (string, error)
	UploadJSON(v interface{}) (string, error)
	Fetch(cid string) ([]byte, error)
}

// ShellClient implements IPFSClient via go-ipfs-api.
type ShellClient struct {
	shell *shell.Shell
}

// NewIPFSShell creates an IPFS shell client targeting the given endpoint.
func NewIPFSShell(endpoint string) (*ShellClient, error) {
	if endpoint == "" {
		return nil, errors.New("ipfs: endpoint required")
	}

	return &ShellClient{shell: shell.NewShell(endpoint)}, nil
}

// UploadBytes writes arbitrary bytes to IPFS and returns the resulting CID.
func (c *ShellClient) UploadBytes(data []byte) (string, error) {
	if c == nil || c.shell == nil {
		return "", errors.New("ipfs: shell not initialised")
	}

	return c.shell.Add(bytes.NewReader(data))
}

// UploadJSON marshals the value to JSON and uploads it to IPFS.
func (c *ShellClient) UploadJSON(v interface{}) (string, error) {
	if c == nil || c.shell == nil {
		return "", errors.New("ipfs: shell not initialised")
	}

	payload, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return c.UploadBytes(payload)
}

// Fetch retrieves data from IPFS by CID.
func (c *ShellClient) Fetch(cid string) ([]byte, error) {
	if c == nil || c.shell == nil {
		return nil, errors.New("ipfs: shell not initialised")
	}

	if cid == "" {
		return nil, errors.New("ipfs: cid required")
	}

	reader, err := c.shell.Cat(cid)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// MemoryIPFS provides an in-memory IPFS implementation useful for tests.
type MemoryIPFS struct {
	mu    sync.RWMutex
	store map[string][]byte
}

// NewMemoryIPFS constructs an in-memory store.
func NewMemoryIPFS() *MemoryIPFS {
	return &MemoryIPFS{store: make(map[string][]byte)}
}

// UploadBytes stores the bytes under a deterministic key and returns it as a pseudo CID.
func (m *MemoryIPFS) UploadBytes(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("ipfs: data required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	key := utils.ComputeSHA256(data)[:16]
	m.store[key] = append([]byte(nil), data...)
	return key, nil
}

// UploadJSON marshals the value then stores it.
func (m *MemoryIPFS) UploadJSON(v interface{}) (string, error) {
	payload, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return m.UploadBytes(payload)
}

// Fetch returns the stored bytes.
func (m *MemoryIPFS) Fetch(cid string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, ok := m.store[cid]
	if !ok {
		return nil, errors.New("ipfs: cid not found")
	}

	return append([]byte(nil), data...), nil
}
