package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/catService/internal/model"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func SetupRedis() RedisRepository {
	var db *redis.Client
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("redis", "3.2", nil)
	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		db = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resource.GetHostPort("6379/tcp")),
		})

		return db.Ping(context.Background()).Err()
	}); err != nil {
		logrus.Fatalf("Could not connect to docker: %s", err)
	}

	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		logrus.Fatalf("Could not purge resource: %s", err)
	}

	newRedisRepository := NewRedisRepository(db)

	return newRedisRepository
}

func TestRedisSet(t *testing.T) {
	catRedisRepository := SetupRedis()

	t.Run("SetCat", func(t *testing.T) {
		catID := uuid.New()
		cat := &model.Cat{
			ID:         catID,
			Name:       "Cat 1",
			Age:        2,
			Vaccinated: true,
		}

		err := catRedisRepository.Create(context.Background(), cat)
		require.NoError(t, err)
		require.Nil(t, err)
	})
}
