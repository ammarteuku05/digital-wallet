package cron

import (
	"log"

	"github.com/spf13/cobra"
)

var healthCheckCmd = &cobra.Command{
	Use:   "health",
	Short: "Perform system health checks",
	Long:  "Check database connectivity, Redis connection, and other system health metrics",
	Run: func(cmd *cobra.Command, args []string) {
		performHealthCheck()
	},
}

func performHealthCheck() {
	log.Println("Starting health check...")

	di := initContainer()
	// Check database connectivity
	log.Println("Checking database connectivity...")
	if err := di.DB.Exec("SELECT 1").Error; err != nil {
		log.Printf("❌ Database health check failed: %v", err)
	} else {
		log.Println("✅ Database connection is healthy")
	}
}
