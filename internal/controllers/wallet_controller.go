package controllers

import (
	"digital-wallet/di"
	"digital-wallet/internal/dto"
	"digital-wallet/internal/interfaces"
	"digital-wallet/pkg/response"
	"fmt"

	"github.com/labstack/echo/v4"
)

type WalletController struct {
	walletService interfaces.WalletService
}

func NewWalletController(di *di.Container) *WalletController {
	return &WalletController{
		walletService: di.WalletService,
	}
}

// GetBalance is
func (wc *WalletController) GetBalance(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("user_id")

	res, err := wc.walletService.GetBalance(ctx, userID)
	if err != nil {
		return response.GenerateResponseFromIError(err)
	}

	return response.OK(c, "Balance retrieved successfully", res)
}

// Withdraw is
func (wc *WalletController) Withdraw(c echo.Context) error {
	var req dto.WithdrawRequest
	ctx := c.Request().Context()

	// Bind request
	if err := c.Bind(&req); err != nil {
		return response.ErrBadRequest(err)
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return response.NewValidationError(err.Error())
	}

	// Process withdrawal
	res, err := wc.walletService.Withdraw(ctx, req)
	if err != nil {
		return response.GenerateResponseFromIError(err)
	}

	return response.OK(c, "Withdrawal processed successfully", res)
}

// GetTransactionHistory is
func (wc *WalletController) GetTransactionHistory(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("user_id")

	limit := 10
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	if o := c.QueryParam("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}

	transactions, total, err := wc.walletService.GetTransactionHistory(ctx, userID, limit, offset)
	if err != nil {
		return response.GenerateResponseFromIError(err)
	}

	// Map to DTO
	txResponses := make([]dto.TransactionHistoryResponse, len(transactions))
	for i, tx := range transactions {
		txResponses[i] = dto.TransactionHistoryResponse{
			ID:          tx.ID,
			Amount:      tx.Amount,
			Type:        tx.Type,
			Status:      tx.Status,
			Description: tx.Description,
			CreatedAt:   tx.CreatedAt.String(),
		}
	}

	res := dto.PaginatedTransactionResponse{
		Data: txResponses,
		Meta: dto.PaginationMeta{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	}

	return response.OK(c, "Transaction history retrieved successfully", res)
}
