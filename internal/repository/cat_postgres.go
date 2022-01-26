// Package repository ...
package repository

import (
	"context"
	"fmt"

	"github.com/catService/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

// CatPostgresRepository contains a link to the connection to db
type CatPostgresRepository struct {
	db *pgxpool.Pool
}

// NewCatPostgres create new instance
func NewCatPostgres(pool *pgxpool.Pool) *CatPostgresRepository {
	return &CatPostgresRepository{db: pool}
}

// Get returns cat
func (r *CatPostgresRepository) Get(ctx context.Context, id uuid.UUID) (*model.Cat, error) {
	cat := model.Cat{}
	row := r.db.QueryRow(ctx, "SELECT * FROM cats WHERE id = $1", id)

	err := row.Scan(&cat.ID, &cat.Name, &cat.Age, &cat.Vaccinated)
	if err != nil {
		return nil, fmt.Errorf("get method error %w", err)
	}

	return &cat, nil
}

// Create new cat in db
func (r *CatPostgresRepository) Create(ctx context.Context, cat *model.Cat) error {
	_, err := r.db.Exec(ctx, "INSERT INTO cats(id, name, age, vaccinated) VALUES ($1,$2,$3,$4)",
		cat.ID, cat.Name, cat.Age, cat.Vaccinated)
	if err != nil {
		return fmt.Errorf("create method error %w", err)
	}

	return nil
}

// Delete cat from db
func (r *CatPostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM cats WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete method error %w", err)
	}

	return nil
}

// Update states for cat
func (r *CatPostgresRepository) Update(ctx context.Context, cat *model.Cat) error {
	_, err := r.db.Exec(ctx, "UPDATE cats SET name=$1, age=$2, vaccinated=$3 WHERE id=$4", cat.Name, cat.Age, cat.Vaccinated, cat.ID)
	if err != nil {
		return fmt.Errorf("update method error %w", err)
	}

	return nil
}
