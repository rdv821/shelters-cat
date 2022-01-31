package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/catService/internal/model"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CatRedisCache struct
type CatRedisCache struct {
	client *redis.Client
	mutex  sync.RWMutex
	cats   map[string]*model.Cat
}

// NewRedisCache ...
func NewRedisCache(ctx context.Context, client *redis.Client) *CatRedisCache {
	cats := make(map[string]*model.Cat)
	cache := CatRedisCache{
		client: client,
		cats:   cats,
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data, err := client.XRead(ctx, &redis.XReadArgs{
					Streams: []string{"cats", "$"},
				}).Result()
				if err != nil {
					logrus.Errorf("XREAD result error %v", err)
				}
				for _, result := range data {
					for _, value := range result.Messages {
						for key, value := range value.Values {
							err := cache.handleAction(key, value)
							if err != nil {
								logrus.Errorf("handle error %v", err)
							}
						}
					}
				}
			}
		}
	}()
	return &cache
}

func (c *CatRedisCache) handleAction(action string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	switch action {
	case "create":
		cat := model.Cat{}
		val, ok := value.(string)
		if !ok {
			return errors.New("cast error")
		}
		err := json.Unmarshal([]byte(val), &cat)
		if err != nil {
			return errors.New("cant unmarshal value")
		}
		c.cats[cat.ID.String()] = &model.Cat{
			ID:         cat.ID,
			Name:       cat.Name,
			Age:        cat.Age,
			Vaccinated: cat.Vaccinated,
		}
	case "delete":
		k, ok := value.(string)
		if !ok {
			return errors.New("cast error")
		}
		if _, exist := c.cats[k]; !exist {
			return errors.New("key not found")
		}
		delete(c.cats, k)
	}

	return nil
}

// Get return cat
func (c *CatRedisCache) Get(id fmt.Stringer) (*model.Cat, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	cat, exist := c.cats[id.String()]
	if !exist {
		return nil, errors.New("cat don't exist")
	}

	return cat, nil
}

// Create add new cat in stream
func (c *CatRedisCache) Create(ctx context.Context, cat *model.Cat) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	marshalCat, err := cat.MarshalBinary()
	if err != nil {
		return err
	}

	err = c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "cats",
		MaxLen: 0,
		Values: map[string]interface{}{
			"create": marshalCat,
		},
	}).Err()

	return err
}

// Delete add id in stream for removing
func (c *CatRedisCache) Delete(ctx context.Context, id uuid.UUID) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "cats",
		MaxLen: 0,
		Values: map[string]interface{}{
			"delete": id,
		},
	}).Err()

	return err
}
