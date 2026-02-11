package repositories

import (
	"context"
	"digital-wallet/internal/interfaces"
	"digital-wallet/internal/models"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestWalletRepository_Create_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletRepository(db)

	t.Run("successfully create wallet", func(t *testing.T) {
		wallet := &models.Wallet{
			ID:       "wallet-1",
			UserID:   "user-1",
			Balance:  1000,
			Currency: "IDR",
			IsActive: true,
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `wallets`")).
			WithArgs(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.Background(), wallet)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create wallet error", func(t *testing.T) {
		wallet := &models.Wallet{ID: "wallet-err"}
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `wallets`")).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Create(context.Background(), wallet)
		assert.Error(t, err)
	})
}

func TestWalletRepository_GetByID_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletRepository(db)

	t.Run("successfully get wallet by id", func(t *testing.T) {
		walletID := "wallet-1"
		rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "is_active"}).
			AddRow(walletID, "user-1", 1000.0, "IDR", true)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallets` WHERE id = ? AND `wallets`.`deleted_at` IS NULL")).
			WithArgs(walletID, 1).
			WillReturnRows(rows)

		wallet, err := repo.GetByID(context.Background(), walletID)
		assert.NoError(t, err)
		assert.NotNil(t, wallet)
		assert.Equal(t, walletID, wallet.ID)
	})

	t.Run("get wallet by id error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallets` WHERE id = ?")).
			WillReturnError(errors.New("db error"))

		wallet, err := repo.GetByID(context.Background(), "err")
		assert.Error(t, err)
		assert.Nil(t, wallet)
	})
}

func TestWalletRepository_GetByUserID_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletRepository(db)

	t.Run("successfully get wallet by user id", func(t *testing.T) {
		userID := "user-1"
		rows := sqlmock.NewRows([]string{"id", "user_id", "balance", "currency", "is_active"}).
			AddRow("wallet-1", userID, 1000.0, "IDR", true)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallets` WHERE user_id = ? AND `wallets`.`deleted_at` IS NULL")).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		wallet, err := repo.GetByUserID(context.Background(), userID)
		assert.NoError(t, err)
		assert.NotNil(t, wallet)
		assert.Equal(t, userID, wallet.UserID)
	})

	t.Run("user wallet not found", func(t *testing.T) {
		userID := "unknown"
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `wallets` WHERE user_id = ? AND `wallets`.`deleted_at` IS NULL")).
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		wallet, err := repo.GetByUserID(context.Background(), userID)
		assert.NoError(t, err)
		assert.Nil(t, wallet)
	})

	t.Run("get by user id error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(errors.New("err"))
		_, err := repo.GetByUserID(context.Background(), "err")
		assert.Error(t, err)
	})
}

func TestWalletRepository_GetBalance_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletRepository(db)

	t.Run("successfully get balance", func(t *testing.T) {
		walletID := "wallet-1"
		rows := sqlmock.NewRows([]string{"balance"}).AddRow(1500.50)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT `balance` FROM `wallets` WHERE id = ? AND `wallets`.`deleted_at` IS NULL")).
			WithArgs(walletID, 1).
			WillReturnRows(rows)

		balance, err := repo.GetBalance(context.Background(), walletID)
		assert.NoError(t, err)
		assert.Equal(t, 1500.50, balance)
	})

	t.Run("get balance error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(errors.New("err"))
		_, err := repo.GetBalance(context.Background(), "err")
		assert.Error(t, err)
	})
}

func TestWalletRepository_UpdateBalance_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletRepository(db)

	t.Run("successfully update balance", func(t *testing.T) {
		walletID := "wallet-1"
		amount := 100.0

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `wallets` SET").
			WithArgs(amount, sqlmock.AnyArg(), walletID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.UpdateBalance(context.Background(), walletID, amount)
		assert.NoError(t, err)
	})

	t.Run("update balance error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE").WillReturnError(errors.New("err"))
		mock.ExpectRollback()

		err := repo.UpdateBalance(context.Background(), "err", 10)
		assert.Error(t, err)
	})
}

func TestWalletRepository_Update_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletRepository(db)

	t.Run("successfully update wallet", func(t *testing.T) {
		wallet := &models.Wallet{ID: "wallet-1", Balance: 2000}
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `wallets` SET").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Update(context.Background(), wallet)
		assert.NoError(t, err)
	})
}

func TestWalletRepository_Withdraw_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletRepository(db)

	t.Run("successfully withdraw within transaction", func(t *testing.T) {
		walletID := "wallet-1"
		amount := 500.0

		mock.ExpectBegin()
		mock.ExpectQuery("SELECT \\* FROM `wallets` WHERE id = \\?").
			WithArgs(walletID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "balance", "is_active"}).
				AddRow(walletID, "user-1", 1000.0, true))

		mock.ExpectExec("UPDATE `wallets` SET").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectQuery("SELECT \\* FROM `wallets` WHERE id = \\?").
			WithArgs(walletID, walletID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "balance"}).AddRow(walletID, 500.0))
		mock.ExpectCommit()

		wallet, err := repo.Withdraw(context.Background(), walletID, amount)
		assert.NoError(t, err)
		assert.NotNil(t, wallet)
		assert.Equal(t, 500.0, wallet.Balance)
	})

	t.Run("wallet not found in withdraw", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT").WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		_, err := repo.Withdraw(context.Background(), "none", 100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wallet not found")
	})

	t.Run("inactive wallet in withdraw", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "is_active"}).AddRow("w1", false))
		mock.ExpectRollback()

		_, err := repo.Withdraw(context.Background(), "w1", 100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wallet is not active")
	})

	t.Run("update fails in withdraw", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "is_active"}).AddRow("w1", 1000.0, true))
		mock.ExpectExec("UPDATE").WillReturnError(errors.New("update err"))
		mock.ExpectRollback()

		_, err := repo.Withdraw(context.Background(), "w1", 100)
		assert.Error(t, err)
	})

	t.Run("rows affected zero in withdraw", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "balance", "is_active"}).AddRow("w1", 1000.0, true))
		mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectRollback()

		_, err := repo.Withdraw(context.Background(), "w1", 100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update wallet balance")
	})
}

func TestWalletTransactionRepository_Create_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletTransactionRepository(db)

	t.Run("successfully create transaction", func(t *testing.T) {
		tx := &models.WalletTransaction{ID: "tx-1"}
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `wallet_transactions`").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.Background(), tx)
		assert.NoError(t, err)
	})
}

func TestWalletTransactionRepository_GetByID_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletTransactionRepository(db)

	t.Run("successfully get by id", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("tx-1"))
		res, err := repo.GetByID(context.Background(), "tx-1")
		assert.NoError(t, err)
		assert.Equal(t, "tx-1", res.ID)
	})

	t.Run("get by id error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(errors.New("err"))
		_, err := repo.GetByID(context.Background(), "tx-1")
		assert.Error(t, err)
	})
}

func TestWalletTransactionRepository_GetByWalletID_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletTransactionRepository(db)

	t.Run("successfully get transactions by wallet id", func(t *testing.T) {
		walletID := "wallet-1"
		rows := sqlmock.NewRows([]string{"id", "wallet_id", "amount"}).
			AddRow("tx-1", walletID, 100.0)

		mock.ExpectQuery("SELECT \\* FROM `wallet_transactions` WHERE wallet_id = \\?").
			WithArgs(walletID, 10).
			WillReturnRows(rows)

		txs, err := repo.GetByWalletID(context.Background(), walletID, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, txs, 1)
	})
}

func TestWalletTransactionRepository_Update_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewWalletTransactionRepository(db)

	t.Run("successfully update transaction", func(t *testing.T) {
		tx := &models.WalletTransaction{ID: "tx-1"}
		mock.ExpectBegin()
		mock.ExpectExec("UPDATE `wallet_transactions` SET").WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Update(context.Background(), tx)
		assert.NoError(t, err)
	})
}

func TestRepositoryRegistry_DoInTransaction_Real(t *testing.T) {
	db, mock := setupMockDB(t)
	registry := NewRepositoryRegistry(db, nil, nil)

	t.Run("successfully execute in transaction", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit()

		out, err := registry.DoInTransaction(context.Background(), func(ctx context.Context, txRepo interfaces.RegistryRepository) (interface{}, error) {
			return "success", nil
		})

		assert.NoError(t, err)
		assert.Equal(t, "success", out)
	})

	t.Run("rollback on panic", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback()

		defer func() {
			recover()
		}()

		_, _ = registry.DoInTransaction(context.Background(), func(ctx context.Context, txRepo interfaces.RegistryRepository) (interface{}, error) {
			panic("oops")
		})
	})
}
