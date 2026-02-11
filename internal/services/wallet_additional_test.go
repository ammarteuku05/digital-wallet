package services

import (
	"context"
	"errors"
	"testing"

	"digital-wallet/configs"
	"digital-wallet/internal/dto"
	"digital-wallet/internal/mocks"
	"digital-wallet/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestWalletService_Withdraw_InactiveWallet tests withdrawal from inactive wallet with actual service call
func TestWalletService_Withdraw_InactiveWallet_Real(t *testing.T) {
	mockWalletRepo := mocks.NewWalletRepository(t)
	mockTxRepo := mocks.NewWalletTransactionRepository(t)

	inactiveWallet := &models.Wallet{
		ID:       "wallet-inactive",
		UserID:   "user-inactive",
		Balance:  1000.00,
		Currency: "IDR",
		IsActive: false, // Inactive wallet
	}

	mockWalletRepo.On("GetByUserID", mock.Anything, "user-inactive").Return(inactiveWallet, nil)

	reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
	svc := NewWalletService(reg, (*configs.Config)(nil))

	req := dto.WithdrawRequest{UserID: "user-inactive", Amount: 100.00, Description: "test withdrawal"}
	resp, err := svc.Withdraw(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "not active")
}

// TestWalletService_Withdraw_TransactionCreateError tests error when transaction creation fails
func TestWalletService_Withdraw_TransactionCreateError(t *testing.T) {
	mockWalletRepo := mocks.NewWalletRepository(t)
	mockTxRepo := mocks.NewWalletTransactionRepository(t)

	wallet := &models.Wallet{
		ID:       "wallet-1",
		UserID:   "user-1",
		Balance:  1000.00,
		Currency: "IDR",
		IsActive: true,
	}

	mockWalletRepo.On("GetByUserID", mock.Anything, "user-1").Return(wallet, nil)
	mockTxRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("transaction create error"))

	reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
	svc := NewWalletService(reg, (*configs.Config)(nil))

	req := dto.WithdrawRequest{UserID: "user-1", Amount: 100.00, Description: "test withdrawal"}
	resp, err := svc.Withdraw(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, resp)
}

// TestWalletService_Withdraw_WithdrawalFails tests when withdrawal operation fails
func TestWalletService_Withdraw_WithdrawalFails(t *testing.T) {
	mockWalletRepo := mocks.NewWalletRepository(t)
	mockTxRepo := mocks.NewWalletTransactionRepository(t)

	wallet := &models.Wallet{
		ID:       "wallet-1",
		UserID:   "user-1",
		Balance:  100.00,
		Currency: "IDR",
		IsActive: true,
	}

	mockWalletRepo.On("GetByUserID", mock.Anything, "user-1").Return(wallet, nil)
	mockTxRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	mockWalletRepo.On("Withdraw", mock.Anything, "wallet-1", 500.00).Return(nil, errors.New("insufficient balance"))
	mockTxRepo.On("Update", mock.Anything, mock.MatchedBy(func(tx *models.WalletTransaction) bool {
		return tx.Status == "FAILED"
	})).Return(nil)

	reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
	svc := NewWalletService(reg, (*configs.Config)(nil))

	req := dto.WithdrawRequest{UserID: "user-1", Amount: 500.00, Description: "test withdrawal"}
	resp, err := svc.Withdraw(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, resp)
	assert.Contains(t, err.Error(), "withdrawal failed")
}

// TestWalletService_Withdraw_TransactionUpdateError tests error when transaction update fails
func TestWalletService_Withdraw_TransactionUpdateError(t *testing.T) {
	mockWalletRepo := mocks.NewWalletRepository(t)
	mockTxRepo := mocks.NewWalletTransactionRepository(t)

	wallet := &models.Wallet{
		ID:       "wallet-1",
		UserID:   "user-1",
		Balance:  1000.00,
		Currency: "IDR",
		IsActive: true,
	}

	updatedWallet := &models.Wallet{
		ID:       "wallet-1",
		UserID:   "user-1",
		Balance:  500.00,
		Currency: "IDR",
		IsActive: true,
	}

	mockWalletRepo.On("GetByUserID", mock.Anything, "user-1").Return(wallet, nil)
	mockTxRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
	mockWalletRepo.On("Withdraw", mock.Anything, "wallet-1", 500.00).Return(updatedWallet, nil)
	mockTxRepo.On("Update", mock.Anything, mock.MatchedBy(func(tx *models.WalletTransaction) bool {
		return tx.Status == "COMPLETED"
	})).Return(errors.New("update error"))

	reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
	svc := NewWalletService(reg, (*configs.Config)(nil))

	req := dto.WithdrawRequest{UserID: "user-1", Amount: 500.00, Description: "test withdrawal"}
	resp, err := svc.Withdraw(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, resp)
}

// TestWalletService_Withdraw_GetOrCreateWalletError tests error in GetOrCreateWallet
func TestWalletService_Withdraw_GetOrCreateWalletError(t *testing.T) {
	mockWalletRepo := mocks.NewWalletRepository(t)
	mockTxRepo := mocks.NewWalletTransactionRepository(t)

	mockWalletRepo.On("GetByUserID", mock.Anything, "user-error").Return(nil, errors.New("database connection error"))

	reg := &testRegistry{wr: mockWalletRepo, tr: mockTxRepo}
	svc := NewWalletService(reg, (*configs.Config)(nil))

	req := dto.WithdrawRequest{UserID: "user-error", Amount: 100.00, Description: "test withdrawal"}
	resp, err := svc.Withdraw(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, resp)
}
