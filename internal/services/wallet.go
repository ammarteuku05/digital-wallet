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

// GetOrCreateWallet is
func (s *WalletService) GetOrCreateWallet(ctx context.Context, userID string) (*models.Wallet, error) {
	walletRepo := s.repo.GetWalletRepository()

	// check if exists
	wallet, err := walletRepo.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, response.Wrap(err, "error retrieving wallet")
	}

	if wallet != nil {
		return wallet, nil
	}

	// create
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

// GetBalance is
func (s *WalletService) GetBalance(ctx context.Context, userID string) (*dto.BalanceResponse, error) {
	walletRepo := s.repo.GetWalletRepository()

	// GetOrCreateWallet is
	wallet, err := s.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, err
	}

	// GetBalance is
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

// Withdraw is
func (s *WalletService) Withdraw(ctx context.Context, req dto.WithdrawRequest) (*dto.WithdrawResponse, error) {
	// GetOrCreateWallet is
	wallet, err := s.GetOrCreateWallet(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// check if active
	if !wallet.IsActive {
		return nil, response.NewValidationError("Wallet is not active")
	}

	result, err := s.repo.DoInTransaction(ctx, func(ctx context.Context, txRepo interfaces.RegistryRepository) (interface{}, error) {
		walletRepo := txRepo.GetWalletRepository()
		transactionRepo := txRepo.GetWalletTransactionRepository()

		// Create
		transactionID := uuid.New().String()
		transaction := &models.WalletTransaction{
			ID:          transactionID,
			WalletID:    wallet.ID,
			Amount:      req.Amount,
			Type:        "WITHDRAWAL",
			Status:      "PENDING",
			Description: req.Description,
		}

		// Create transaction
		if err := transactionRepo.Create(ctx, transaction); err != nil {
			return nil, response.Wrap(err, "error creating transaction")
		}

		// withdrawal with row-level locking
		updatedWallet, err := walletRepo.Withdraw(ctx, wallet.ID, req.Amount)
		if err != nil {
			// if withdrawal fails
			transaction.Status = "FAILED"
			_ = transactionRepo.Update(ctx, transaction)
			return nil, response.Wrap(err, "withdrawal failed")
		}

		// transaction is completed
		transaction.Status = "COMPLETED"
		if err := transactionRepo.Update(ctx, transaction); err != nil {
			return nil, response.Wrap(err, "error updating transaction status")
		}

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

// GetTransactionHistory is
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

	total, err := transactionRepo.CountByWalletID(ctx, wallet.ID)
	if err != nil {
		return nil, 0, response.Wrap(err, "error counting transactions")
	}

	transactions, err := transactionRepo.GetByWalletID(ctx, wallet.ID, limit, offset)
	if err != nil {
		return nil, 0, response.Wrap(err, "error retrieving transactions")
	}

	return transactions, total, nil
}
