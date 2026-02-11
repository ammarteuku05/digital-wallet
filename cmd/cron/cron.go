package cron

import (
	"digital-wallet/di"

	"github.com/spf13/cobra"
)

var CronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Run cron jobs",
	Long:  "Execute scheduled cron jobs for data cleanup, reports, etc.",
}

func init() {
	// Add subcommands for different cron jobs
	CronCmd.AddCommand(healthCheckCmd)
}

// Helper function to initialize di for cron jobs
func initContainer() *di.Container {
	// Initialize di
	return di.SetUp()
}
