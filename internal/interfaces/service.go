package interfaces

import (
	"context"
	"digital-wallet/internal/dto"
	"digital-wallet/internal/models"
)

//go:generate mockery --name WalletService --case snake --output ../mocks --disable-version-string

// WalletService interface
type WalletService interface {
	GetOrCreateWallet(ctx context.Context, userID string) (*models.Wallet, error)
	GetBalance(ctx context.Context, userID string) (*dto.BalanceResponse, error)
	Withdraw(ctx context.Context, req dto.WithdrawRequest) (*dto.WithdrawResponse, error)
	GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]models.WalletTransaction, int64, error)
}
