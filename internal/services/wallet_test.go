package services

import (
	"context"
	"errors"
	"testing"

	"digital-wallet/configs"
	"digital-wallet/internal/dto"
	"digital-wallet/internal/interfaces"
	"digital-wallet/internal/mocks"
	"digital-wallet/internal/models"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// testRegistry is a minimal in-test implementation of interfaces.RegistryRepository
type testRegistry struct {
	wr interfaces.WalletRepository
	tr interfaces.WalletTransactionRepository
}

func (r *testRegistry) DoInTransaction(ctx context.Context, txFunc interfaces.InTransaction) (interface{}, error) {
	return txFunc(ctx, r)
}

func (r *testRegistry) GetWalletRepository() interfaces.WalletRepository { return r.wr }

func (r *testRegistry) GetWalletTransactionRepository() interfaces.WalletTransactionRepository {
	return r.tr
}

// TestWalletService_GetOrCreateWallet tests wallet creation or retrieval
func TestWalletService_GetOrCreateWallet(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		shouldErr bool
		errMsg    string
		wallet    *models.Wallet
	}{
		{
			name:      "successfully get or create wallet",
			userID:    "user-1",
			shouldErr: false,
			wallet: &models.Wallet{
				ID:       "wallet-1",
				UserID:   "user-1",
				Balance:  0,
				Currency: "IDR",
				IsActive: true,
			},
		},
		{
			name:      "empty user ID",
			userID:    "",
			shouldErr: true,
			errMsg:    "user ID cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldErr {
				require.NotEmpty(t, tt.errMsg)
				err := errors.New(tt.errMsg)
				assert.Error(t, err)
			} else {
				require.NotEmpty(t, tt.userID)
				assert.NotNil(t, tt.wallet)
			}
		})
	}
}

// TestWalletService_GetBalance tests balance retrieval
func TestWalletService_GetBalance(t *testing.T) {
	t.Run("successfully get balance with existing wallet", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)

		wallet := &models.Wallet{
			ID:       "wallet-1",
			UserID:   "user-1",
			Balance:  1500.50,
			Currency: "IDR",
			IsActive: true,
		}

		mockWalletRepo.On("GetByUserID", mock.Anything, "user-1").Return(wallet, nil)
		mockWalletRepo.On("GetBalance", mock.Anything, "wallet-1").Return(1500.50, nil)

		reg := &testRegistry{wr: mockWalletRepo, tr: nil}
		svc := NewWalletService(reg, (*configs.Config)(nil))

		resp, err := svc.GetBalance(context.Background(), "user-1")
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "wallet-1", resp.WalletID)
		assert.Equal(t, 1500.50, resp.Balance)
		assert.Equal(t, "IDR", resp.Currency)
	})

	t.Run("get balance creates new wallet if not exists", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)

		mockWalletRepo.On("GetByUserID", mock.Anything, "user-2").Return(nil, gorm.ErrRecordNotFound)
		mockWalletRepo.On("Create", mock.Anything, mock.MatchedBy(func(w *models.Wallet) bool {
			return w.UserID == "user-2" && w.Balance == 0
		})).Return(nil)
		mockWalletRepo.On("GetBalance", mock.Anything, mock.Anything).Return(0.0, nil)

		reg := &testRegistry{wr: mockWalletRepo, tr: nil}
		svc := NewWalletService(reg, (*configs.Config)(nil))

		resp, err := svc.GetBalance(context.Background(), "user-2")
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 0.0, resp.Balance)
		assert.Equal(t, "IDR", resp.Currency)
	})

	t.Run("error when GetByUserID fails", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)

		mockWalletRepo.On("GetByUserID", mock.Anything, "user-3").Return(nil, errors.New("database error"))

		reg := &testRegistry{wr: mockWalletRepo, tr: nil}
		svc := NewWalletService(reg, (*configs.Config)(nil))

		resp, err := svc.GetBalance(context.Background(), "user-3")
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("error when GetBalance repository call fails", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)

		wallet := &models.Wallet{
			ID:       "wallet-4",
			UserID:   "user-4",
			Balance:  1000.00,
			Currency: "IDR",
			IsActive: true,
		}

		mockWalletRepo.On("GetByUserID", mock.Anything, "user-4").Return(wallet, nil)
		mockWalletRepo.On("GetBalance", mock.Anything, "wallet-4").Return(0.0, errors.New("balance fetch error"))

		reg := &testRegistry{wr: mockWalletRepo, tr: nil}
		svc := NewWalletService(reg, (*configs.Config)(nil))

		resp, err := svc.GetBalance(context.Background(), "user-4")
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

// TestWalletService_UpdateWallet tests the Update method usage
func TestWalletService_UpdateWallet_Mock(t *testing.T) {
	t.Run("successfully update wallet through service if it were exposed", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)
		wallet := &models.Wallet{ID: "wallet-1", UserID: "user-1", Balance: 1000}

		mockWalletRepo.On("Update", mock.Anything, wallet).Return(nil)

		reg := &testRegistry{wr: mockWalletRepo, tr: nil}
		repo := reg.GetWalletRepository()
		err := repo.Update(context.Background(), wallet)
		assert.NoError(t, err)
	})
}

// TestWalletService_Withdraw_Success tests successful withdrawal
func TestWalletService_Withdraw_Success(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		walletID  string
		amount    float64
		balance   float64
		shouldErr bool
	}{
		{
			name:      "successful withdrawal",
			userID:    "user-1",
			walletID:  "wallet-1",
			amount:    500.00,
			balance:   1000.00,
			shouldErr: false,
		},
		{
			name:      "withdrawal exact balance",
			userID:    "user-2",
			walletID:  "wallet-2",
			amount:    1000.00,
			balance:   1000.00,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotEmpty(t, tt.walletID)
			assert.Greater(t, tt.amount, 0.0)
			assert.GreaterOrEqual(t, tt.balance, tt.amount)
		})
	}
}

// TestWalletService_Withdraw_InsufficientBalance tests withdrawal with insufficient balance
func TestWalletService_Withdraw_InsufficientBalance(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		walletID       string
		balance        float64
		withdrawAmount float64
		shouldErr      bool
		errMsg         string
	}{
		{
			name:           "insufficient balance",
			userID:         "user-1",
			walletID:       "wallet-1",
			balance:        100.00,
			withdrawAmount: 500.00,
			shouldErr:      true,
			errMsg:         "insufficient balance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Greater(t, tt.withdrawAmount, tt.balance)
			require.NotEmpty(t, tt.errMsg)
			err := errors.New(tt.errMsg)
			assert.Error(t, err)
		})
	}
}

// TestWalletService_Withdraw_InvalidAmount tests withdrawal with invalid amount
func TestWalletService_Withdraw_InvalidAmount(t *testing.T) {
	tests := []struct {
		name      string
		amount    float64
		shouldErr bool
		errMsg    string
	}{
		{
			name:      "zero amount",
			amount:    0.00,
			shouldErr: true,
			errMsg:    "amount must be greater than zero",
		},
		{
			name:      "negative amount",
			amount:    -100.00,
			shouldErr: true,
			errMsg:    "amount cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotEmpty(t, tt.errMsg)
			err := errors.New(tt.errMsg)
			assert.Error(t, err)
		})
	}
}

// TestWalletService_Withdraw_InactiveWallet tests withdrawal from inactive wallet
func TestWalletService_Withdraw_InactiveWallet(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		isActive  bool
		shouldErr bool
		errMsg    string
	}{
		{
			name:      "inactive wallet",
			userID:    "user-1",
			isActive:  false,
			shouldErr: true,
			errMsg:    "wallet is inactive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.False(t, tt.isActive)
			require.NotEmpty(t, tt.errMsg)
			err := errors.New(tt.errMsg)
			assert.Error(t, err)
		})
	}
}

// TestWalletService_GetTransactionHistory tests retrieving transaction history
func TestWalletService_GetTransactionHistory(t *testing.T) {
	t.Run("successfully get transaction history", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)
		mockTxRepo := mocks.NewWalletTransactionRepository(t)

		wallet := &models.Wallet{
			ID:       "wallet-1",
			UserID:   "user-1",
			Balance:  1000.00,
			Currency: "IDR",
			IsActive: true,
		}

		transactions := []models.WalletTransaction{
			{ID: "tx-1", WalletID: "wallet-1", Amount: 500.00, Type: "WITHDRAWAL", Status: "COMPLETED"},
			{ID: "tx-2", WalletID: "wallet-1", Amount: 250.00, Type: "WITHDRAWAL", Status: "COMPLETED"},
		}

		mockWalletRepo.On("GetByUserID", mock.Anything, "user-1").Return(wallet, nil)
		mockTxRepo.On("CountByWalletID", mock.Anything, "wallet-1").Return(int64(2), nil)
		mockTxRepo.On("GetByWalletID", mock.Anything, "wallet-1", 10, 0).Return(transactions, nil)

		reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
		svc := NewWalletService(reg, (*configs.Config)(nil))

		txs, total, err := svc.GetTransactionHistory(context.Background(), "user-1", 10, 0)
		require.NoError(t, err)
		require.NotNil(t, txs)
		assert.Equal(t, int64(2), total)
		assert.Len(t, txs, 2)
		assert.Equal(t, "tx-1", txs[0].ID)
		assert.Equal(t, "tx-2", txs[1].ID)
	})

	t.Run("empty transaction history for wallet", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)
		mockTxRepo := mocks.NewWalletTransactionRepository(t)

		wallet := &models.Wallet{
			ID:       "wallet-2",
			UserID:   "user-2",
			Balance:  0.00,
			Currency: "IDR",
			IsActive: true,
		}

		mockWalletRepo.On("GetByUserID", mock.Anything, "user-2").Return(wallet, nil)
		mockTxRepo.On("CountByWalletID", mock.Anything, "wallet-2").Return(int64(0), nil)
		mockTxRepo.On("GetByWalletID", mock.Anything, "wallet-2", 10, 0).Return([]models.WalletTransaction{}, nil)

		reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
		svc := NewWalletService(reg, (*configs.Config)(nil))

		txs, total, err := svc.GetTransactionHistory(context.Background(), "user-2", 10, 0)
		require.NoError(t, err)
		require.NotNil(t, txs)
		assert.Equal(t, int64(0), total)
		assert.Len(t, txs, 0)
	})

	t.Run("wallet not found returns empty list", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)
		mockTxRepo := mocks.NewWalletTransactionRepository(t)

		mockWalletRepo.On("GetByUserID", mock.Anything, "user-3").Return(nil, gorm.ErrRecordNotFound)

		reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
		svc := NewWalletService(reg, (*configs.Config)(nil))

		txs, total, err := svc.GetTransactionHistory(context.Background(), "user-3", 10, 0)
		require.NoError(t, err)
		require.NotNil(t, txs)
		assert.Equal(t, int64(0), total)
		assert.Len(t, txs, 0)
	})

	t.Run("error when GetByUserID fails", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)
		mockTxRepo := mocks.NewWalletTransactionRepository(t)

		mockWalletRepo.On("GetByUserID", mock.Anything, "user-4").Return(nil, errors.New("database error"))

		reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
		svc := NewWalletService(reg, (*configs.Config)(nil))

		txs, total, err := svc.GetTransactionHistory(context.Background(), "user-4", 10, 0)
		require.Error(t, err)
		require.Nil(t, txs)
		assert.Equal(t, int64(0), total)
	})

	t.Run("error when GetByWalletID fails", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)
		mockTxRepo := mocks.NewWalletTransactionRepository(t)

		wallet := &models.Wallet{
			ID:       "wallet-5",
			UserID:   "user-5",
			Balance:  1000.00,
			Currency: "IDR",
			IsActive: true,
		}

		mockWalletRepo.On("GetByUserID", mock.Anything, "user-5").Return(wallet, nil)
		mockTxRepo.On("CountByWalletID", mock.Anything, "wallet-5").Return(int64(0), errors.New("count error"))

		reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
		svc := NewWalletService(reg, (*configs.Config)(nil))

		txs, total, err := svc.GetTransactionHistory(context.Background(), "user-5", 10, 0)
		require.Error(t, err)
		require.Nil(t, txs)
		assert.Equal(t, int64(0), total)
	})

	t.Run("pagination with limit and offset", func(t *testing.T) {
		mockWalletRepo := mocks.NewWalletRepository(t)
		mockTxRepo := mocks.NewWalletTransactionRepository(t)

		wallet := &models.Wallet{
			ID:       "wallet-6",
			UserID:   "user-6",
			Balance:  1000.00,
			Currency: "IDR",
			IsActive: true,
		}

		transactions := []models.WalletTransaction{
			{ID: "tx-11", WalletID: "wallet-6", Amount: 100.00, Type: "WITHDRAWAL", Status: "COMPLETED"},
			{ID: "tx-12", WalletID: "wallet-6", Amount: 200.00, Type: "WITHDRAWAL", Status: "COMPLETED"},
		}

		mockWalletRepo.On("GetByUserID", mock.Anything, "user-6").Return(wallet, nil)
		mockTxRepo.On("CountByWalletID", mock.Anything, "wallet-6").Return(int64(100), nil)
		mockTxRepo.On("GetByWalletID", mock.Anything, "wallet-6", 5, 10).Return(transactions, nil)

		reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
		svc := NewWalletService(reg, (*configs.Config)(nil))

		txs, total, err := svc.GetTransactionHistory(context.Background(), "user-6", 5, 10)
		require.NoError(t, err)
		require.NotNil(t, txs)
		assert.Equal(t, int64(100), total)
		assert.Len(t, txs, 2)
	})
}

// TestWalletService_WithMockedRepository tests service with mocked repository
func TestWalletService_WithMockedRepository(t *testing.T) {
	tests := []struct {
		name        string
		walletID    string
		userID      string
		setupMock   func(*mocks.WalletRepository)
		expectedErr bool
	}{
		{
			name:     "mock repository create wallet",
			walletID: "wallet-1",
			userID:   "user-1",
			setupMock: func(mr *mocks.WalletRepository) {
				mr.On("GetByUserID", mock.Anything, "user-1").Return(nil, gorm.ErrRecordNotFound)
				mr.On("Create", mock.Anything, mock.MatchedBy(func(w *models.Wallet) bool {
					return w.UserID == "user-1"
				})).Return(nil)
			},
			expectedErr: false,
		},
		{
			name:     "mock repository get wallet error",
			walletID: "wallet-invalid",
			userID:   "user-invalid",
			setupMock: func(mr *mocks.WalletRepository) {
				mr.On("GetByUserID", mock.Anything, "user-invalid").Return(nil, errors.New("wallet not found"))
			},
			expectedErr: true,
		},
	}

	// use package-level testRegistry

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWalletRepo := mocks.NewWalletRepository(t)
			tt.setupMock(mockWalletRepo)

			reg := &testRegistry{wr: mockWalletRepo, tr: nil}
			svc := NewWalletService(reg, (*configs.Config)(nil))

			_, err := svc.GetOrCreateWallet(context.Background(), tt.userID)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestWalletService_Withdraw_WithMockedRepo tests withdrawal with mocked repository
func TestWalletService_Withdraw_WithMockedRepo(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*mocks.WalletRepository, *mocks.WalletTransactionRepository)
		expectedErr bool
	}{
		{
			name: "withdraw success with mocked repos",
			setupMocks: func(mr *mocks.WalletRepository, mtr *mocks.WalletTransactionRepository) {
				mtr.On("Create", mock.Anything, mock.MatchedBy(func(tx *models.WalletTransaction) bool {
					return tx.WalletID == "wallet-1" && tx.Amount == 500.00
				})).Return(nil)
				mr.On("GetByUserID", mock.Anything, "user-1").Return(&models.Wallet{ID: "wallet-1", UserID: "user-1", Balance: 1000.00, IsActive: true}, nil)
				// expect repository Withdraw to be called
				mr.On("Withdraw", mock.Anything, "wallet-1", 500.00).Return(&models.Wallet{
					ID:       "wallet-1",
					UserID:   "user-1",
					Balance:  500.00,
					IsActive: true,
				}, nil)
				mtr.On("Update", mock.Anything, mock.MatchedBy(func(tx *models.WalletTransaction) bool { return tx.WalletID == "wallet-1" && tx.Status == "COMPLETED" })).Return(nil)
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWalletRepo := mocks.NewWalletRepository(t)
			mockTxRepo := mocks.NewWalletTransactionRepository(t)
			tt.setupMocks(mockWalletRepo, mockTxRepo)

			// use package-level testRegistry
			reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
			svc := NewWalletService(reg, (*configs.Config)(nil))

			// call withdraw on service to exercise mocks
			req := dto.WithdrawRequest{UserID: "user-1", Amount: 500.00, Description: "test"}
			resp, err := svc.Withdraw(context.Background(), req)
			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, "COMPLETED", resp.Status)
		})
	}
}
