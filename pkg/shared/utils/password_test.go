package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashAndSalt(t *testing.T) {
	t.Run("successfully hash password", func(t *testing.T) {
		password := []byte("mySecurePassword123")
		hash, err := HashAndSalt(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, string(password), hash)
	})

	t.Run("hash empty password", func(t *testing.T) {
		password := []byte("")
		hash, err := HashAndSalt(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("hash long password (>72 chars)", func(t *testing.T) {
		// bcrypt has a 72 character limit
		longPassword := []byte("thisIsAVeryLongPasswordThatExceedsTheSeventyTwoCharacterLimitForBcryptHashingAlgorithm")
		hash, err := HashAndSalt(longPassword)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("same password produces different hashes", func(t *testing.T) {
		password := []byte("samePassword")
		hash1, err1 := HashAndSalt(password)
		hash2, err2 := HashAndSalt(password)

		require.NoError(t, err1)
		require.NoError(t, err2)
		// Due to salt, same password should produce different hashes
		assert.NotEqual(t, hash1, hash2)
	})
}

func TestComparePasswords(t *testing.T) {
	t.Run("correct password matches hash", func(t *testing.T) {
		password := []byte("correctPassword123")
		hash, err := HashAndSalt(password)
		require.NoError(t, err)

		result := ComparePasswords(hash, password)
		assert.True(t, result)
	})

	t.Run("incorrect password does not match hash", func(t *testing.T) {
		password := []byte("correctPassword123")
		wrongPassword := []byte("wrongPassword456")
		hash, err := HashAndSalt(password)
		require.NoError(t, err)

		result := ComparePasswords(hash, wrongPassword)
		assert.False(t, result)
	})

	t.Run("empty password comparison", func(t *testing.T) {
		password := []byte("password")
		hash, err := HashAndSalt(password)
		require.NoError(t, err)

		result := ComparePasswords(hash, []byte(""))
		assert.False(t, result)
	})

	t.Run("invalid hash returns false", func(t *testing.T) {
		invalidHash := "not-a-valid-bcrypt-hash"
		password := []byte("password")

		result := ComparePasswords(invalidHash, password)
		assert.False(t, result)
	})
}
