package services

import (
	"context"
	"digital-wallet/configs"
	"digital-wallet/internal/dto"
	"digital-wallet/internal/interfaces"
	"digital-wallet/internal/models"
	response "digital-wallet/pkg/response"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletService struct {
	repo interfaces.RegistryRepository
	cfg  *configs.Config
}

// Ensure WalletService implements interfaces.WalletService
var _ interfaces.WalletService = (*WalletService)(nil)

func NewWalletService(repo interfaces.RegistryRepository, config *configs.Config) interfaces.WalletService {
	return &WalletService{
		repo: repo,
		cfg:  config,
	}
}

// GetOrCreateWallet retrieves or creates a wallet for a user
func (s *WalletService) GetOrCreateWallet(ctx context.Context, userID string) (*models.Wallet, error) {
	walletRepo := s.repo.GetWalletRepository()

	// Check if wallet already exists
	wallet, err := walletRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, response.Wrap(err, "error retrieving wallet")
	}

	// If wallet exists, return it
	if wallet != nil {
		return wallet, nil
	}

	// Create new wallet
	newWallet := &models.Wallet{
		ID:       uuid.New().String(),
		UserID:   userID,
		Balance:  0,
		Currency: "IDR",
		IsActive: true,
	}

	if err := walletRepo.Create(ctx, newWallet); err != nil {
		return nil, response.Wrap(err, "error creating wallet")
	}

	return newWallet, nil
}

// GetBalance returns the current balance of a user's wallet
func (s *WalletService) GetBalance(ctx context.Context, userID string) (*dto.BalanceResponse, error) {
	walletRepo := s.repo.GetWalletRepository()

	// Get or create wallet
	wallet, err := s.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get balance
	balance, err := walletRepo.GetBalance(ctx, wallet.ID)
	if err != nil {
		return nil, response.Wrap(err, "error retrieving balance")
	}

	return &dto.BalanceResponse{
		WalletID: wallet.ID,
		Balance:  balance,
		Currency: wallet.Currency,
		IsActive: wallet.IsActive,
	}, nil
}

// Withdraw processes a withdrawal from a wallet using concurrency-safe mechanism
// It leverages database transactions and row-level locking to ensure thread-safety
func (s *WalletService) Withdraw(ctx context.Context, req dto.WithdrawRequest) (*dto.WithdrawResponse, error) {
	// Get or create wallet
	wallet, err := s.GetOrCreateWallet(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// Check if wallet is active
	if !wallet.IsActive {
		return nil, response.NewValidationError("Wallet is not active")
	}

	// Perform withdrawal within a database transaction
	result, err := s.repo.DoInTransaction(ctx, func(ctx context.Context, txRepo interfaces.RegistryRepository) (interface{}, error) {
		walletRepo := txRepo.GetWalletRepository()
		transactionRepo := txRepo.GetWalletTransactionRepository()

		// Create transaction record (mark as pending)
		transactionID := uuid.New().String()
		transaction := &models.WalletTransaction{
			ID:          transactionID,
			WalletID:    wallet.ID,
			Amount:      req.Amount,
			Type:        "WITHDRAWAL",
			Status:      "PENDING",
			Description: req.Description,
		}

		// Create transaction in database
		if err := transactionRepo.Create(ctx, transaction); err != nil {
			return nil, response.Wrap(err, "error creating transaction")
		}

		// Perform concurrency-safe withdrawal with row-level locking
		updatedWallet, err := walletRepo.Withdraw(ctx, wallet.ID, req.Amount)
		if err != nil {
			// Mark transaction as failed if withdrawal fails
			transaction.Status = "FAILED"
			_ = transactionRepo.Update(ctx, transaction)
			return nil, response.Wrap(err, "withdrawal failed")
		}

		// Mark transaction as completed
		transaction.Status = "COMPLETED"
		if err := transactionRepo.Update(ctx, transaction); err != nil {
			return nil, response.Wrap(err, "error updating transaction status")
		}

		// Return response data
		return &dto.WithdrawResponse{
			ID:            wallet.ID,
			WalletID:      wallet.ID,
			Amount:        req.Amount,
			NewBalance:    updatedWallet.Balance,
			TransactionID: transactionID,
			Status:        "COMPLETED",
			Timestamp:     transaction.UpdatedAt.String(),
		}, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*dto.WithdrawResponse), nil
}

// GetTransactionHistory returns transaction history for a wallet
func (s *WalletService) GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]models.WalletTransaction, int64, error) {
	walletRepo := s.repo.GetWalletRepository()
	transactionRepo := s.repo.GetWalletTransactionRepository()

	// Get wallet
	wallet, err := walletRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, response.Wrap(err, "error retrieving wallet")
	}

	if wallet == nil {
		return []models.WalletTransaction{}, 0, nil
	}

	// Get total count
	total, err := transactionRepo.CountByWalletID(ctx, wallet.ID)
	if err != nil {
		return nil, 0, response.Wrap(err, "error counting transactions")
	}

	// Get transactions
	transactions, err := transactionRepo.GetByWalletID(ctx, wallet.ID, limit, offset)
	if err != nil {
		return nil, 0, response.Wrap(err, "error retrieving transactions")
	}

	return transactions, total, nil
}
