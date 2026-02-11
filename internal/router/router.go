package router

import (
	"digital-wallet/di"
	"digital-wallet/internal/controllers"

	"github.com/labstack/echo/v4"
)

func SetupRouter(e *echo.Echo, di *di.Container) {
	// Initialize controllers
	walletController := controllers.NewWalletController(di)

	v1 := e.Group("/v1")
	{

		// Wallet routes
		wallet := v1.Group("/wallet")
		// wallet.Use(middleware.AuthMiddleware(di))
		{
			wallet.GET("/balance/:user_id", walletController.GetBalance)
			wallet.POST("/withdraw", walletController.Withdraw)
			wallet.GET("/:user_id/transactions", walletController.GetTransactionHistory)

		}

		e.Any("", func(c echo.Context) error {
			return echo.NotFoundHandler(c)
		})

		e.Any("/*", func(c echo.Context) error {
			return echo.NotFoundHandler(c)
		})
	}
}
