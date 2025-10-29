package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

const nonceSize = 12

func deriveAESKey(passphrase string) []byte {
	sum := sha256.Sum256([]byte(passphrase))
	key := make([]byte, len(sum))
	copy(key, sum[:])
	return key[:32]
}

func encryptString(aesKey []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("wallet: cipher init failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("wallet: gcm init failed: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("wallet: nonce generation failed: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	payload := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(payload), nil
}

func decryptString(aesKey []byte, payload string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", fmt.Errorf("wallet: decode failed: %w", err)
	}

	if len(data) < nonceSize {
		return "", fmt.Errorf("wallet: payload too short")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("wallet: cipher init failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("wallet: gcm init failed: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("wallet: decrypt failed: %w", err)
	}

	return string(plaintext), nil
}
