package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/catService/internal/model"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// CatRedisRepository contains a link to the connection to db
type CatRedisRepository struct {
	client *redis.Client
}

// NewCatRedis create new instance
func NewCatRedis(client *redis.Client) *CatRedisRepository {
	return &CatRedisRepository{client: client}
}

// Get returns cat
func (r *CatRedisRepository) Get(ctx context.Context, id uuid.UUID) (*model.Cat, error) {
	cat := model.Cat{}

	redisCat, err := r.client.Get(ctx, id.String()).Result()
	if err == redis.Nil || err != nil {
		return nil, fmt.Errorf("get redis error: %w", err)
	}

	err = json.Unmarshal([]byte(redisCat), &cat)
	if err != nil {
		return nil, err
	}

	return &cat, nil
}

// Create new cat in redis
func (r *CatRedisRepository) Create(ctx context.Context, cat *model.Cat) error {
	jsonCat, err := json.Marshal(cat)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, cat.ID.String(), jsonCat, 0).Err()
	if err != nil {
		return fmt.Errorf("create redis error: %w", err)
	}
	return nil
}

// Delete cat from db
func (r *CatRedisRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.Del(ctx, id.String()).Err()
	if err != nil {
		return fmt.Errorf("delete redis error: %w", err)
	}
	return nil
}
