package router

import (
	"digital-wallet/di"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	e := echo.New()
	container := &di.Container{}

	SetupRouter(e, container)

	t.Run("Verify route registrations", func(t *testing.T) {
		routes := e.Routes()
		foundBalance := false
		foundWithdraw := false
		foundHistory := false

		for _, r := range routes {
			if r.Path == "/v1/wallet/balance/:user_id" && r.Method == http.MethodGet {
				foundBalance = true
			}
			if r.Path == "/v1/wallet/withdraw" && r.Method == http.MethodPost {
				foundWithdraw = true
			}
			if r.Path == "/v1/wallet/:user_id/transactions" && r.Method == http.MethodGet {
				foundHistory = true
			}
		}

		assert.True(t, foundBalance, "Balance route not found")
		assert.True(t, foundWithdraw, "Withdraw route not found")
		assert.True(t, foundHistory, "History route not found")
	})
}
