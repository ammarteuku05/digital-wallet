package main

import (
	"digital-wallet/cmd/api"
	"digital-wallet/cmd/cron"
	"digital-wallet/cmd/migrate"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "digital-wallet",
		Short: "digital-wallet Test Application",
		Long:  "A comprehensive guest experience management system with API server and cron job capabilities",
	}

	// Add subcommands
	rootCmd.AddCommand(api.ServerCmd)
	rootCmd.AddCommand(cron.CronCmd)
	rootCmd.AddCommand(migrate.MigrateCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
