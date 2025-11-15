package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

type CacheService struct {
	client *redis.Client
}

func NewCacheService() (*CacheService, error) {
	cacheUrl, ok := os.LookupEnv("CACHE_URL")
	if !ok {
		return nil, fmt.Errorf("CACHE_URL is not set")
	}

	opt, err := redis.ParseURL(cacheUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_URL: %w", err)
	}

	// BUG: This is an issue from redis/go-redis. Update once fixed
	// https://github.com/redis/go-redis/issues/3536
	opt.MaintNotificationsConfig = &maintnotifications.Config{
		Mode: maintnotifications.ModeDisabled,
	}

	return &CacheService{client: redis.NewClient(opt)}, nil
}

func (c *CacheService) SetOTP(ctx context.Context, id uuid.UUID, otp string) error {
	err := c.client.Set(ctx, id.String(), otp, 30*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheService) GetOTP(ctx context.Context, id uuid.UUID) (string, error) {
	otp, err := c.client.Get(ctx, id.String()).Result()
	if err != nil {
		return "", err
	}

	return otp, nil
}
