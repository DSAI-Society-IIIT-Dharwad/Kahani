package storage_test

import (
	"testing"

	"storytelling-blockchain/internal/storage"
)

func TestMemoryIPFS(t *testing.T) {
	ipfs := storage.NewMemoryIPFS()

	cid, err := ipfs.UploadBytes([]byte("hello"))
	if err != nil {
		t.Fatalf("upload bytes failed: %v", err)
	}

	jsonCID, err := ipfs.UploadJSON(map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatalf("upload json failed: %v", err)
	}

	if cid == jsonCID {
		t.Fatalf("expected unique cids for distinct payloads")
	}

	data, err := ipfs.Fetch(cid)
	if err != nil {
		t.Fatalf("fetch failed: %v", err)
	}

	if string(data) != "hello" {
		t.Fatalf("unexpected payload: %s", string(data))
	}

	if _, err := ipfs.Fetch("missing"); err == nil {
		t.Fatalf("expected error for missing cid")
	}
}

func TestNewIPFSShellValidation(t *testing.T) {
	if _, err := storage.NewIPFSShell(""); err == nil {
		t.Fatalf("expected error when endpoint empty")
	}
}
