package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/yellow78/mini-mes/backend/internal/model"
)

// RecipeRepository Recipe DB 存取介面
type RecipeRepository interface {
	FindAll(ctx context.Context) ([]model.Recipe, error)
}

type recipeRepo struct {
	db *sqlx.DB
}

// NewRecipeRepository 建立 Recipe repository 實例
func NewRecipeRepository(db *sqlx.DB) RecipeRepository {
	return &recipeRepo{db: db}
}

// FindAll 取得所有 Recipe
func (r *recipeRepo) FindAll(ctx context.Context) ([]model.Recipe, error) {
	const q = `
		SELECT id, name, equipment_type, target_temp, target_pressure, duration_min, created_at
		FROM recipe
		ORDER BY equipment_type, name
	`
	var recipes []model.Recipe
	err := r.db.SelectContext(ctx, &recipes, q)
	return recipes, err
}
