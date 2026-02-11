package api

import (
	"digital-wallet/configs"
	"digital-wallet/di"
	"digital-wallet/internal/router"
	"digital-wallet/pkg/response"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

var (
	port string
)

var ServerCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server",
	Long:  "Start the digital-wallet API server with Echo framework",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	ServerCmd.Flags().StringVarP(&port, "port", "p", "", "Port to run the server on (overrides config)")
}

func startServer() {
	// Load configuration
	cfg := configs.LoadDefault()

	// Override port if provided via flag
	if port != "" {
		cfg.Server.PORT = port
	}

	// Initialize di
	di := di.SetUp()

	// Initialize Echo
	e := echo.New()

	// Set custom error handler
	e.HTTPErrorHandler = response.CustomHTTPErrorHandler

	// validation
	e.Validator = di.Validator

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(response.EchoMiddleware())

	// Setup routes
	router.SetupRouter(e, di)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.PORT)
	if err := e.Start(":" + cfg.Server.PORT); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
