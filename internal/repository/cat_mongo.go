// Package repository ...
package repository

import (
	"context"
	"fmt"

	"github.com/catService/internal/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CatMongoRepository contains a link to the connection to db
type CatMongoRepository struct {
	db *mongo.Database
}

// NewCatMongo create new instance
func NewCatMongo(collection *mongo.Database) *CatMongoRepository {
	return &CatMongoRepository{db: collection}
}

// Get returns cat
func (c *CatMongoRepository) Get(ctx context.Context, id uuid.UUID) (*model.Cat, error) {
	cat := model.Cat{}
	result := c.db.Collection("cat").FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return nil, fmt.Errorf("get method error %w", result.Err())
	}

	if err := result.Decode(&cat); err != nil {
		return nil, fmt.Errorf("failed decode cat from DB %w", err)
	}

	return &cat, nil
}

// Create new cat in db
func (c *CatMongoRepository) Create(ctx context.Context, cat *model.Cat) error {
	_, err := c.db.Collection("cat").InsertOne(ctx, &cat)
	if err != nil {
		return fmt.Errorf("create method error %w", err)
	}

	return nil
}

// Update states for cat
func (c *CatMongoRepository) Update(ctx context.Context, cat *model.Cat) error {
	filter := bson.M{"_id": cat.ID}

	update := bson.M{"name": cat.Name, "age": cat.Age, "vaccinated": cat.Vaccinated}

	_, err := c.db.Collection("cat").UpdateOne(ctx, filter, bson.M{
		"$set": update,
	})
	if err != nil {
		return fmt.Errorf("failed to execute update cat query: %w", err)
	}

	return nil
}

// Delete cat from db
func (c *CatMongoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	filter := bson.M{"_id": id}

	_, err := c.db.Collection("cat").DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("delete method error %w", err)
	}

	return nil
}
