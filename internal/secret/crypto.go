package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

const encryptedPrefix = "kron:v1:"

type Manager interface {
	Encrypt(value string) (string, error)
	Decrypt(value string) (string, error)
	IsEncrypted(value string) bool
}

type NoopManager struct{}

func (NoopManager) Encrypt(value string) (string, error) { return value, nil }
func (NoopManager) Decrypt(value string) (string, error) { return value, nil }
func (NoopManager) IsEncrypted(value string) bool        { return false }

type AESGCMManager struct {
	gcm cipher.AEAD
}

func NewAESGCMManager(keyValue string) (*AESGCMManager, error) {
	keyValue = strings.TrimSpace(keyValue)
	if keyValue == "" {
		return nil, errors.New("encryption key is required")
	}

	key, err := base64.StdEncoding.DecodeString(keyValue)
	if err != nil || len(key) != 32 {
		key = []byte(keyValue)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes or base64-encoded 32 bytes, got %d bytes", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &AESGCMManager{gcm: gcm}, nil
}

func (m *AESGCMManager) Encrypt(value string) (string, error) {
	if m.IsEncrypted(value) {
		return value, nil
	}

	nonce := make([]byte, m.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := m.gcm.Seal(nil, nonce, []byte(value), nil)
	payload := append(nonce, ciphertext...)
	return encryptedPrefix + base64.StdEncoding.EncodeToString(payload), nil
}

func (m *AESGCMManager) Decrypt(value string) (string, error) {
	if !m.IsEncrypted(value) {
		return value, nil
	}

	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, encryptedPrefix))
	if err != nil {
		return "", err
	}

	nonceSize := m.gcm.NonceSize()
	if len(payload) < nonceSize {
		return "", errors.New("encrypted value is too short")
	}

	nonce := payload[:nonceSize]
	ciphertext := payload[nonceSize:]
	plaintext, err := m.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (m *AESGCMManager) IsEncrypted(value string) bool {
	return strings.HasPrefix(value, encryptedPrefix)
}
