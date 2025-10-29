package utils

import "testing"

func TestEd25519RoundTrip(t *testing.T) {
	pub, priv, err := GenerateEd25519Keypair()
	if err != nil {
		t.Fatalf("failed to generate keypair: %v", err)
	}

	sig, err := SignEd25519(priv, []byte("payload"))
	if err != nil {
		t.Fatalf("failed to sign payload: %v", err)
	}

	valid, err := VerifyEd25519(pub, []byte("payload"), sig)
	if err != nil {
		t.Fatalf("verification failed: %v", err)
	}

	if !valid {
		t.Fatalf("expected signature to verify")
	}

	invalid, err := VerifyEd25519(pub, []byte("payload"), "not-a-sig")
	if err == nil {
		t.Fatalf("expected error when decoding invalid signature, got validity %v", invalid)
	}
}
