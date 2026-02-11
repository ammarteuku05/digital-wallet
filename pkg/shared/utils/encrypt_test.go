package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	t.Run("encode bytes to base64", func(t *testing.T) {
		input := []byte("hello world")
		encoded := Encode(input)

		assert.NotEmpty(t, encoded)
		assert.Equal(t, "aGVsbG8gd29ybGQ=", encoded)
	})

	t.Run("encode empty bytes", func(t *testing.T) {
		input := []byte("")
		encoded := Encode(input)

		assert.Equal(t, "", encoded)
	})
}

func TestEncrypt(t *testing.T) {
	t.Run("successfully encrypt text", func(t *testing.T) {
		text := "sensitive data"
		// AES-128 requires 16 byte key, AES-256 requires 32 byte key
		secret := "1234567890123456" // 16 bytes

		encrypted, err := Encrypt(text, secret)

		require.NoError(t, err)
		assert.NotEmpty(t, encrypted)
		assert.NotEqual(t, text, encrypted)
	})

	t.Run("encrypt empty text", func(t *testing.T) {
		text := ""
		secret := "1234567890123456"

		encrypted, err := Encrypt(text, secret)

		require.NoError(t, err)
		assert.Empty(t, encrypted)
	})

	t.Run("same text produces same encrypted output with same key", func(t *testing.T) {
		text := "test data"
		secret := "1234567890123456"

		encrypted1, err1 := Encrypt(text, secret)
		encrypted2, err2 := Encrypt(text, secret)

		require.NoError(t, err1)
		require.NoError(t, err2)
		// With CFB mode and same IV (bytess), same input produces same output
		assert.Equal(t, encrypted1, encrypted2)
	})

	t.Run("error with invalid key length", func(t *testing.T) {
		text := "test"
		invalidSecret := "short" // Invalid key length

		encrypted, err := Encrypt(text, invalidSecret)

		require.Error(t, err)
		assert.Empty(t, encrypted)
	})

	t.Run("encrypt with 32 byte key (AES-256)", func(t *testing.T) {
		text := "secret message"
		secret := "12345678901234567890123456789012" // 32 bytes

		encrypted, err := Encrypt(text, secret)

		require.NoError(t, err)
		assert.NotEmpty(t, encrypted)
	})
}
