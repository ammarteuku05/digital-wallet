package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Set required env vars
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_BASE_URL", "http://localhost")
	os.Setenv("APP_NAME", "test-app")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USERNAME", "root")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "test_db")
	os.Setenv("DB_MIN_IDDLE_CONN", "1")
	os.Setenv("DB_MAX_OPEN_CONN", "10")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("REDIS_DB", "0")
	os.Setenv("JWT_SIGNING_KEY", "secret")
	os.Setenv("JWT_TOKEN_EXPIRATION", "3600")
	os.Setenv("JWT_REFRESH_TOKEN_EXPIRATION_DAY", "7")
	os.Setenv("JWT_ENCRYPTION_KEY", "12345678901234567890123456789012")
	os.Setenv("ACCOUNT_ACTIVATION_TOKEN_EXPIRATION", "24h")
	os.Setenv("FORGOT_PASSWORD_TOKEN_EXPIRATION", "1")
	os.Setenv("LOGGER_STDOUT", "true")
	os.Setenv("LOGGER_FILE_LOCATION", "./logs")
	os.Setenv("LOGGER_FILE_MAX_AGE", "7")
	os.Setenv("LOGGER_LEVEL", "2")
	os.Setenv("LOGGER_MASKING", "true")
	os.Setenv("LOGGER_MASKING_PARAMS", "password,token")
	t.Run("LoadDefault", func(t *testing.T) {
		cfg := LoadDefault()
		assert.NotNil(t, cfg)
		assert.Equal(t, Env("test"), cfg.Server.ENV)
	})

	t.Run("LoadTest", func(t *testing.T) {
		cfg := LoadTest()
		assert.NotNil(t, cfg)
		assert.Equal(t, "test-app", cfg.Server.NAME)
	})
}
