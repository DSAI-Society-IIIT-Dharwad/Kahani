package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

// GenerateEd25519Keypair yields a base64 encoded keypair.
func GenerateEd25519Keypair() (publicKey string, privateKey string, err error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(pub), base64.StdEncoding.EncodeToString(priv), nil
}

// SignEd25519 signs the payload using the provided base64 encoded private key.
func SignEd25519(privateKey string, payload []byte) (string, error) {
	privBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	if len(privBytes) != ed25519.PrivateKeySize {
		return "", errors.New("invalid ed25519 private key length")
	}

	signature := ed25519.Sign(ed25519.PrivateKey(privBytes), payload)
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyEd25519 validates the payload signature using the base64 encoded public key.
func VerifyEd25519(publicKey string, payload []byte, signature string) (bool, error) {
	pubBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false, err
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	if len(pubBytes) != ed25519.PublicKeySize {
		return false, errors.New("invalid ed25519 public key length")
	}

	return ed25519.Verify(ed25519.PublicKey(pubBytes), payload, sigBytes), nil
}
