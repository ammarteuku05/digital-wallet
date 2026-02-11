package di

import (
	"digital-wallet/configs"
	"fmt"
	"log/slog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MySQLConn(cfg *configs.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true&multiStatements=true",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		slog.Error("Failed to get database instance", "error", err)
		panic(err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(cfg.Database.MinIdleConnections)

	return conn
}
