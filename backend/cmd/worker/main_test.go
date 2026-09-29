package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/winnerx0/kron/internal/domain"
)

var errReadFailed = errors.New("read failed")

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errReadFailed }

func TestReadResponseBody_LimitsBodyTo64KiB(t *testing.T) {
	body := strings.Repeat("x", 64*1024+100)

	got := readResponseBody(strings.NewReader(body))

	require.Len(t, got, 64*1024)
	require.Equal(t, body[:64*1024], got)
}

func TestReadResponseBody_ReadErrorReturnsEmptyString(t *testing.T) {
	require.Empty(t, readResponseBody(failingReader{}))
}

func TestDecryptHeaders_DecryptsSensitiveHeadersAndPreservesNormalHeaders(t *testing.T) {
	input := map[string]any{
		"Authorization": "Bearer token",
		"Content-Type":  "application/json",
		"X-Retry":       3,
	}

	got, err := decryptHeaders(input)

	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"Authorization": "Bearer token",
		"Content-Type":  "application/json",
		"X-Retry":       "3",
	}, got)
}

func TestDecryptJobSecrets_ReturnsJobWithStringHeaders(t *testing.T) {
	input := domain.Job{ID: "job-1", Headers: map[string]any{
		"Content-Type": "application/json",
		"X-Number":     42,
	}}

	got, err := decryptJobSecrets(input)

	require.NoError(t, err)
	require.Equal(t, input.ID, got.ID)
	require.Equal(t, map[string]any{
		"Content-Type": "application/json",
		"X-Number":     "42",
	}, map[string]any(got.Headers))
}

func TestEncryptJobSecrets_LeavesHeadersUnchangedWithNoopManager(t *testing.T) {
	input := domain.Job{Headers: map[string]any{
		"Authorization": "Bearer token",
		"Content-Type":  "application/json",
	}}

	got, err := encryptJobSecrets(input, nil)

	require.NoError(t, err)
	require.Equal(t, input.Headers, got.Headers)
}
