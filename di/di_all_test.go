package di

import (
	"digital-wallet/configs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	cfg := &configs.Config{}

	t.Run("Test different log levels", func(t *testing.T) {
		levels := []int8{0, 1, 2, 3, 99}
		for _, l := range levels {
			cfg.Logger.Level = l
			logger := newLogger(cfg)
			assert.NotNil(t, logger)
		}
	})
}

func TestSetupCache(t *testing.T) {
	cfg := &configs.Config{}
	cfg.Redis.Host = "localhost"
	cfg.Redis.Port = "6379"

	t.Run("Test cache setup", func(t *testing.T) {
		client := SetupCache(cfg)
		assert.NotNil(t, client)
	})
}
