package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/yellow78/mini-mes/backend/internal/model"
)

type lotRepo struct {
	db *sqlx.DB
}

// NewLotRepository 建立 Lot repository 實例
func NewLotRepository(db *sqlx.DB) LotRepository {
	return &lotRepo{db: db}
}

// FindAll 取得所有 Lot，JOIN recipe 取得 recipe_name
func (r *lotRepo) FindAll(ctx context.Context) ([]model.Lot, error) {
	const q = `
		SELECT l.id, l.lot_number, l.product, l.recipe_id,
		       rc.name AS recipe_name,
		       l.priority, l.status, l.wafer_count, l.created_at
		FROM lot l
		LEFT JOIN recipe rc ON l.recipe_id = rc.id
		ORDER BY l.priority, l.created_at
	`
	rows, err := r.db.QueryxContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Lot
	for rows.Next() {
		var l model.Lot
		var recipeName sql.NullString
		err := rows.Scan(
			&l.ID, &l.LotNumber, &l.Product, &l.RecipeID,
			&recipeName, &l.Priority, &l.Status, &l.WaferCount, &l.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if recipeName.Valid {
			l.RecipeName = &recipeName.String
		}
		result = append(result, l)
	}
	return result, rows.Err()
}

// FindByID 依 id 取得單筆 Lot
func (r *lotRepo) FindByID(ctx context.Context, id int) (*model.Lot, error) {
	const q = `
		SELECT l.id, l.lot_number, l.product, l.recipe_id,
		       rc.name AS recipe_name,
		       l.priority, l.status, l.wafer_count, l.created_at
		FROM lot l
		LEFT JOIN recipe rc ON l.recipe_id = rc.id
		WHERE l.id = $1
	`
	row := r.db.QueryRowxContext(ctx, q, id)
	var l model.Lot
	var recipeName sql.NullString
	err := row.Scan(
		&l.ID, &l.LotNumber, &l.Product, &l.RecipeID,
		&recipeName, &l.Priority, &l.Status, &l.WaferCount, &l.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if recipeName.Valid {
		l.RecipeName = &recipeName.String
	}
	return &l, nil
}

// Create 建立新 Lot，並回傳含 ID 的完整物件
func (r *lotRepo) Create(ctx context.Context, lot *model.Lot) error {
	const q = `
		INSERT INTO lot (lot_number, product, recipe_id, priority, status, wafer_count)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRowxContext(ctx, q,
		lot.LotNumber, lot.Product, lot.RecipeID, lot.Priority, lot.Status, lot.WaferCount,
	).Scan(&lot.ID, &lot.CreatedAt)
}

// UpdateStatus 更新 Lot 狀態
func (r *lotRepo) UpdateStatus(ctx context.Context, id int, status model.LotStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE lot SET status = $1 WHERE id = $2`,
		status, id,
	)
	return err
}
