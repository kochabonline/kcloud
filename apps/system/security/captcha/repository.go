package captcha

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	mobilePrefix = "captcha:mobile:%s"
	emailPrefix  = "captcha:email:%s"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (repo *Repository) key(prefix, key string) string {
	return fmt.Sprintf(prefix, key)
}

func (repo *Repository) GetMobileCode(ctx context.Context, mobile string) (string, time.Duration, error) {
	key := repo.key(mobilePrefix, mobile)

	// 使用 pipeline 获取验证码和剩余时间
	pipe := repo.client.Pipeline()
	getCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return "", 0, err
	}

	code, err := getCmd.Result()
	if err != nil {
		return "", 0, err
	}

	ttl, err := ttlCmd.Result()
	if err != nil {
		return "", 0, err
	}

	return code, ttl, nil
}

func (repo *Repository) SetMobileCode(ctx context.Context, mobile string, code string, expiration time.Duration) error {
	return repo.client.Set(ctx, repo.key(mobilePrefix, mobile), code, 0).Err()
}

func (repo *Repository) GetEmailCode(ctx context.Context, email string) (string, time.Duration, error) {
	key := repo.key(emailPrefix, email)

	// 使用 pipeline 获取验证码和剩余时间
	pipe := repo.client.Pipeline()
	getCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return "", 0, err
	}

	code, err := getCmd.Result()
	if err != nil {
		return "", 0, err
	}

	ttl, err := ttlCmd.Result()
	if err != nil {
		return "", 0, err
	}

	return code, ttl, nil
}

func (repo *Repository) SetEmailCode(ctx context.Context, email string, code string, expiration time.Duration) error {
	return repo.client.Set(ctx, repo.key(emailPrefix, email), code, 0).Err()
}
