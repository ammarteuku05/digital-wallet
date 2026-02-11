package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModels_TableName(t *testing.T) {
	t.Run("Wallet TableName", func(t *testing.T) {
		w := Wallet{}
		assert.Equal(t, "wallets", w.TableName())
	})

	t.Run("WalletTransaction TableName", func(t *testing.T) {
		wt := WalletTransaction{}
		assert.Equal(t, "wallet_transactions", wt.TableName())
	})
}

func TestUser_Methods(t *testing.T) {
	t.Run("HashAndSalt", func(t *testing.T) {
		pwd := []byte("password123")
		hash, err := HashAndSalt(pwd)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	})

	t.Run("CheckPassword", func(t *testing.T) {
		pwd := "password123"
		hash, _ := HashAndSalt([]byte(pwd))
		u := User{Password: hash}

		err := u.CheckPassword(pwd)
		assert.NoError(t, err)

		err = u.CheckPassword("wrong")
		assert.Error(t, err)
	})
}
