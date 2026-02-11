package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"digital-wallet/di"
	"digital-wallet/internal/dto"
	"digital-wallet/internal/mocks"
	"digital-wallet/internal/models"
	"digital-wallet/pkg/response"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWalletController_GetBalance(t *testing.T) {
	e := echo.New()

	t.Run("successfully get balance", func(t *testing.T) {
		mockSvc := mocks.NewWalletService(t)
		container := &di.Container{WalletService: mockSvc, Validator: di.NewCustomValidator()}
		wc := NewWalletController(container)
		e.Validator = container.Validator

		userID := "user-123"
		mockSvc.On("GetBalance", mock.Anything, userID).Return(&dto.BalanceResponse{
			Balance:  1000.0,
			Currency: "IDR",
			WalletID: "wallet-1",
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/balance", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("user_id")
		c.SetParamValues(userID)

		err := wc.GetBalance(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "Balance retrieved successfully", resp["message"])

		data := resp["data"].(map[string]interface{})
		assert.Equal(t, 1000.0, data["balance"])
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := mocks.NewWalletService(t)
		container := &di.Container{WalletService: mockSvc}
		wc := NewWalletController(container)

		userID := "user-456"
		mockSvc.On("GetBalance", mock.Anything, userID).Return(nil, errors.New("something went wrong"))

		req := httptest.NewRequest(http.MethodGet, "/balance", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("user_id")
		c.SetParamValues(userID)

		err := wc.GetBalance(c)
		require.Error(t, err)
		assert.IsType(t, response.ErrorResponse{}, err)
	})
}

func TestWalletController_Withdraw(t *testing.T) {
	e := echo.New()

	t.Run("successfully process withdrawal", func(t *testing.T) {
		mockSvc := mocks.NewWalletService(t)
		container := &di.Container{WalletService: mockSvc, Validator: di.NewCustomValidator()}
		wc := NewWalletController(container)
		e.Validator = container.Validator

		userID := "user-123"
		reqBody := `{"user_id": "user-123", "amount": 500.0, "description": "test"}`
		mockSvc.On("Withdraw", mock.Anything, mock.MatchedBy(func(r dto.WithdrawRequest) bool {
			return r.UserID == userID && r.Amount == 500.0
		})).Return(&dto.WithdrawResponse{
			TransactionID: "tx-1",
			Status:        "COMPLETED",
			NewBalance:    500.0,
		}, nil)

		req := httptest.NewRequest(http.MethodPost, "/withdraw", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userID", userID)

		err := wc.Withdraw(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("bind error", func(t *testing.T) {
		mockSvc := mocks.NewWalletService(t)
		container := &di.Container{WalletService: mockSvc}
		wc := NewWalletController(container)

		reqBody := `invalid json`
		req := httptest.NewRequest(http.MethodPost, "/withdraw", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("userID", "user-1")

		err := wc.Withdraw(c)
		require.Error(t, err)
	})
}

func TestWalletController_GetTransactionHistory(t *testing.T) {
	e := echo.New()

	t.Run("successfully get history", func(t *testing.T) {
		mockSvc := mocks.NewWalletService(t)
		container := &di.Container{WalletService: mockSvc, Validator: di.NewCustomValidator()}
		wc := NewWalletController(container)
		e.Validator = container.Validator

		userID := "user-123"
		mockSvc.On("GetTransactionHistory", mock.Anything, userID, 10, 0).Return([]models.WalletTransaction{
			{ID: "tx-1", Amount: 100.0},
		}, int64(1), nil)

		req := httptest.NewRequest(http.MethodGet, "/v1/wallet/"+userID+"/transactions", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/wallet/:user_id/transactions")
		c.SetParamNames("user_id")
		c.SetParamValues(userID)

		err := wc.GetTransactionHistory(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "\"total\":1")
		assert.Contains(t, rec.Body.String(), "\"limit\":10")
		assert.Contains(t, rec.Body.String(), "\"meta\"")
		assert.Contains(t, rec.Body.String(), "\"offset\":0")
	})

	t.Run("history with custom query params", func(t *testing.T) {
		mockSvc := mocks.NewWalletService(t)
		container := &di.Container{WalletService: mockSvc}
		wc := NewWalletController(container)

		userID := "user-123"
		mockSvc.On("GetTransactionHistory", mock.Anything, userID, 5, 20).Return([]models.WalletTransaction{}, int64(10), nil)

		req := httptest.NewRequest(http.MethodGet, "/v1/wallet/"+userID+"/transactions?limit=5&offset=20", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/wallet/:user_id/transactions")
		c.SetParamNames("user_id")
		c.SetParamValues(userID)

		err := wc.GetTransactionHistory(c)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "\"total\":10")
		assert.Contains(t, rec.Body.String(), "\"limit\":5")
		assert.Contains(t, rec.Body.String(), "\"offset\":20")
	})
}
