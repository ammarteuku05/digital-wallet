package migrate

import (
	"digital-wallet/di"
	"fmt"
	"log"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration(migrate.Up)
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigration(migrate.Down)
	},
}

func init() {
	MigrateCmd.AddCommand(upCmd)
	MigrateCmd.AddCommand(downCmd)
}

func runMigration(direction migrate.MigrationDirection) {
	container := di.SetUp()

	db, err := container.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "scripts",
	}

	n, err := migrate.Exec(db, "mysql", migrations, direction)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
}
