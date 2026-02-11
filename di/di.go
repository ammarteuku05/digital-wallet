package di

import (
	"digital-wallet/configs"
	"digital-wallet/internal/interfaces"
	"digital-wallet/internal/repositories"
	"digital-wallet/internal/services"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	DB            *gorm.DB
	RedisClient   *redis.Client
	Config        *configs.Config
	RepoRegistry  interfaces.RegistryRepository
	WalletService interfaces.WalletService
	Validator     *CustomValidator
	Logger        *slog.Logger
}

func SetUp() *Container {
	var (
		cfg       = configs.LoadDefault()
		validator = NewCustomValidator()
		logger    = newLogger(cfg)
	)

	// initial cache and database
	var (
		redisClient = SetupCache(cfg)
		db          = MySQLConn(cfg)
	)

	// Initialize repositories
	repoRegistry := repositories.NewRepositoryRegistry(db, redisClient, cfg)

	// Initialize services
	walletService := services.NewWalletService(repoRegistry, cfg)

	return &Container{
		DB:            db,
		RedisClient:   redisClient,
		Config:        cfg,
		RepoRegistry:  repoRegistry,
		WalletService: walletService,
		Validator:     validator,
		Logger:        logger,
	}
}
