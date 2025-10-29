package utils

import "testing"

func TestComputeSHA256(t *testing.T) {
	got := ComputeSHA256([]byte("hello"))
	const expected = "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	if got != expected {
		t.Fatalf("unexpected hash: %s", got)
	}
}
