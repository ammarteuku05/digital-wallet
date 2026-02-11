package interfaces

import (
	"context"
	"digital-wallet/internal/models"
)

//go:generate mockery --name WalletRepository --case snake --output ../mocks --disable-version-string

// WalletRepository interface
type WalletRepository interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	GetByID(ctx context.Context, id string) (*models.Wallet, error)
	GetByUserID(ctx context.Context, userID string) (*models.Wallet, error)
	GetBalance(ctx context.Context, walletID string) (float64, error)
	UpdateBalance(ctx context.Context, walletID string, amount float64) error
	Update(ctx context.Context, wallet *models.Wallet) error
	Withdraw(ctx context.Context, walletID string, amount float64) (*models.Wallet, error)
}

//go:generate mockery --name WalletTransactionRepository --case snake --output ../mocks --disable-version-string

// WalletTransactionRepository interface
type WalletTransactionRepository interface {
	Create(ctx context.Context, transaction *models.WalletTransaction) error
	GetByID(ctx context.Context, id string) (*models.WalletTransaction, error)
	GetByWalletID(ctx context.Context, walletID string, limit, offset int) ([]models.WalletTransaction, error)
	CountByWalletID(ctx context.Context, walletID string) (int64, error)
	Update(ctx context.Context, transaction *models.WalletTransaction) error
}

type InTransaction func(ctx context.Context, repoRegistry RegistryRepository) (interface{}, error)
type RegistryRepository interface {
	DoInTransaction(ctx context.Context, txFunc InTransaction) (out interface{}, err error)
	GetWalletRepository() WalletRepository
	GetWalletTransactionRepository() WalletTransactionRepository
}
