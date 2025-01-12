package google

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository struct {
	db       *gorm.DB
	rediscli *redis.Client
}

func NewRepository(db *gorm.DB, rediscli *redis.Client) *Repository {
	return &Repository{
		db:       db,
		rediscli: rediscli,
	}
}

func (repo *Repository) FindByAccountId(ctx context.Context, accountId int64) (*GoogleAuth, error) {
	var auth GoogleAuth
	err := repo.db.WithContext(ctx).Where("account_id = ?", accountId).First(&auth).Error
	return &auth, err
}

func (repo *Repository) Create(ctx context.Context, auth GoogleAuth) error {
	return repo.db.WithContext(ctx).Create(&auth).Error
}

func (repo *Repository) Get(ctx context.Context, key string) (string, error) {
	return repo.rediscli.Get(ctx, key).Result()
}

func (repo *Repository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return repo.rediscli.Set(ctx, key, value, expiration).Err()
}

func (repo *Repository) Del(ctx context.Context, key string) error {
	return repo.rediscli.Del(ctx, key).Err()
}
