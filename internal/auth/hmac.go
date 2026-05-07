package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
)

type HashAlgorithm string

const (
	SHA256 HashAlgorithm = "sha256"
	SHA512 HashAlgorithm = "sha512"
)

// Hash создаёт HMAC-хэш данных с секретным ключом.
// Используй SHA256 для большинства задач, SHA512 — если нужна повышенная стойкость.
func Hash(data, secretKey []byte, alg HashAlgorithm) (string, error) {
	var h func() hash.Hash

	switch alg {
	case SHA256:
		h = sha256.New
	case SHA512:
		h = sha512.New
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", alg)
	}

	mac := hmac.New(h, secretKey)
	if _, err := mac.Write(data); err != nil {
		return "", fmt.Errorf("failed to write data to hmac: %w", err)
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}

// Verify проверяет HMAC-хэш через константное сравнение (защита от timing attack).
func Verify(data, secretKey []byte, expectedHash string, alg HashAlgorithm) (bool, error) {
	actualHash, err := Hash(data, secretKey, alg)
	if err != nil {
		return false, err
	}

	expected, err := hex.DecodeString(expectedHash)
	if err != nil {
		return false, fmt.Errorf("invalid expected hash format: %w", err)
	}

	actual, err := hex.DecodeString(actualHash)
	if err != nil {
		return false, fmt.Errorf("invalid actual hash format: %w", err)
	}

	return hmac.Equal(actual, expected), nil
}
