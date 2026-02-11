package repositories

import (
	"context"
	"digital-wallet/configs"
	"digital-wallet/internal/interfaces"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RepositoryRegistry struct {
	db         *gorm.DB
	redisCache *redis.Client
	cfg        *configs.Config
}

func NewRepositoryRegistry(db *gorm.DB, redisCache *redis.Client, cfg *configs.Config) interfaces.RegistryRepository {
	repo := RepositoryRegistry{
		db,
		redisCache,
		cfg,
	}

	return &repo
}

func (r *RepositoryRegistry) DoInTransaction(ctx context.Context, txFunc interfaces.InTransaction) (out interface{}, err error) {
	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			rErr := tx.Rollback() // err is non-nil; don't change it
			if rErr.Error != nil {
				err = rErr.Error
			}
		} else {
			err = tx.Commit().Error // err is nil; if Commit returns error update err
		}
	}()

	// Create a new registry with the transaction
	txRegistry := &RepositoryRegistry{
		db:         tx,
		redisCache: r.redisCache,
		cfg:        r.cfg,
	}

	// Execute the function with the transaction registry
	out, err = txFunc(ctx, txRegistry)

	return
}

func (r *RepositoryRegistry) GetWalletRepository() interfaces.WalletRepository {
	return NewWalletRepository(r.db)
}

func (r *RepositoryRegistry) GetWalletTransactionRepository() interfaces.WalletTransactionRepository {
	return NewWalletTransactionRepository(r.db)
}
