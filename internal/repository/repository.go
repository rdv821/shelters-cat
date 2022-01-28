// Package repository ...
package repository

import (
	"context"

	"github.com/catService/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

// SheltersCatRepository contains needed methods which must be implemented
//go:generate mockery --dir . --name CatRepository --output ./mocks
type SheltersCatRepository interface {
	Get(context.Context, uuid.UUID) (*model.Cat, error)
	Create(context.Context, *model.Cat) error
	Update(context.Context, *model.Cat) error
	Delete(context.Context, uuid.UUID) error
}

// RedisRepository interface
type RedisRepository interface {
	Get(context.Context, uuid.UUID) (*model.Cat, error)
	Create(context.Context, *model.Cat) error
	Delete(context.Context, uuid.UUID) error
}

// NewPostgresRepository constructor
func NewPostgresRepository(pool *pgxpool.Pool) SheltersCatRepository {
	return NewCatPostgres(pool)
}

// NewMongoRepository constructor
func NewMongoRepository(database *mongo.Database) SheltersCatRepository {
	return NewCatMongo(database)
}

// NewLocalCache constructor
func NewLocalCache(ctx context.Context, client *redis.Client) *CatRedisCache {
	return NewRedisCache(ctx, client)
}
