package repositories

import (
	"context"
	"digital-wallet/internal/interfaces"
	"digital-wallet/internal/models"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository struct {
	db *gorm.DB
}

// Ensure WalletRepository implements interfaces.WalletRepository
var _ interfaces.WalletRepository = (*WalletRepository)(nil)

func NewWalletRepository(database *gorm.DB) interfaces.WalletRepository {
	return &WalletRepository{db: database}
}

func (r *WalletRepository) Create(ctx context.Context, wallet *models.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}

func (r *WalletRepository) GetByID(ctx context.Context, id string) (*models.Wallet, error) {
	var wallet models.Wallet
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&wallet)
	if result.Error != nil {
		return nil, result.Error
	}
	return &wallet, nil
}

func (r *WalletRepository) GetByUserID(ctx context.Context, userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &wallet, nil
}

func (r *WalletRepository) GetBalance(ctx context.Context, walletID string) (float64, error) {
	var wallet models.Wallet
	result := r.db.WithContext(ctx).Select("balance").Where("id = ?", walletID).First(&wallet)
	if result.Error != nil {
		return 0, result.Error
	}
	return wallet.Balance, nil
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, walletID string, amount float64) error {
	return r.db.WithContext(ctx).Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (r *WalletRepository) Update(ctx context.Context, wallet *models.Wallet) error {
	return r.db.WithContext(ctx).Save(wallet).Error
}

func (r *WalletRepository) Withdraw(ctx context.Context, walletID string, amount float64) (*models.Wallet, error) {
	var wallet models.Wallet
	result := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock the wallet row for update (pessimistic locking - FOR UPDATE)
		lockResult := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", walletID).
			First(&wallet)

		if lockResult.Error != nil {
			if errors.Is(lockResult.Error, gorm.ErrRecordNotFound) {
				return errors.New("wallet not found")
			}
			return lockResult.Error
		}

		if !wallet.IsActive {
			return errors.New("wallet is not active")
		}

		if wallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		updateResult := tx.Model(&wallet).Update("balance", gorm.Expr("balance - ?", amount))
		if updateResult.Error != nil {
			return updateResult.Error
		}

		if updateResult.RowsAffected == 0 {
			return errors.New("failed to update wallet balance")
		}

		return tx.First(&wallet, "id = ?", walletID).Error
	})

	if result != nil {
		return nil, result
	}

	return &wallet, nil
}

// WalletTransactionRepository implementation
type WalletTransactionRepository struct {
	db *gorm.DB
}

// Ensure WalletTransactionRepository implements interfaces.WalletTransactionRepository
var _ interfaces.WalletTransactionRepository = (*WalletTransactionRepository)(nil)

func NewWalletTransactionRepository(database *gorm.DB) interfaces.WalletTransactionRepository {
	return &WalletTransactionRepository{db: database}
}

func (r *WalletTransactionRepository) Create(ctx context.Context, transaction *models.WalletTransaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *WalletTransactionRepository) GetByID(ctx context.Context, id string) (*models.WalletTransaction, error) {
	var transaction models.WalletTransaction
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

func (r *WalletTransactionRepository) GetByWalletID(ctx context.Context, walletID string, limit, offset int) ([]models.WalletTransaction, error) {
	var transactions []models.WalletTransaction
	result := r.db.WithContext(ctx).
		Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions)
	return transactions, result.Error
}

func (r *WalletTransactionRepository) CountByWalletID(ctx context.Context, walletID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.WalletTransaction{}).Where("wallet_id = ?", walletID).Count(&count).Error
	return count, err
}

func (r *WalletTransactionRepository) Update(ctx context.Context, transaction *models.WalletTransaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}
