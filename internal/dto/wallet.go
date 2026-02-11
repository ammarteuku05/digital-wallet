package dto

type WithdrawRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	UserID      string  `json:"user_id" validate:"required"`
	Description string  `json:"description"`
}

type WithdrawResponse struct {
	ID            string  `json:"id"`
	WalletID      string  `json:"wallet_id"`
	Amount        float64 `json:"amount"`
	NewBalance    float64 `json:"new_balance"`
	TransactionID string  `json:"transaction_id"`
	Status        string  `json:"status"`
	Timestamp     string  `json:"timestamp"`
}

type BalanceResponse struct {
	WalletID string  `json:"wallet_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	IsActive bool    `json:"is_active"`
}

type TransactionHistoryResponse struct {
	ID          string  `json:"id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

type PaginationMeta struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

type PaginatedTransactionResponse struct {
	Data []TransactionHistoryResponse `json:"data"`
	Meta PaginationMeta               `json:"meta"`
}
