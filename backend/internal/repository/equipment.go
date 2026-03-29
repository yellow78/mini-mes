package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yellow78/mini-mes/backend/internal/model"
)

type equipmentRepo struct {
	db *sqlx.DB
}

// NewEquipmentRepository 建立設備 repository 實例
func NewEquipmentRepository(db *sqlx.DB) EquipmentRepository {
	return &equipmentRepo{db: db}
}

// FindAll 取得所有設備，JOIN lot 取得 lot_number 和 recipe_name
func (r *equipmentRepo) FindAll(ctx context.Context) ([]model.Equipment, error) {
	const q = `
		SELECT
			e.id, e.name, e.type, e.status, e.current_lot_id,
			l.lot_number  AS current_lot,
			rc.name       AS recipe_name,
			e.utilization, e.temperature, e.pressure,
			e.ucl_temp, e.lcl_temp, e.ucl_pressure, e.lcl_pressure,
			e.is_alarm, e.updated_at
		FROM equipment e
		LEFT JOIN lot l    ON e.current_lot_id = l.id
		LEFT JOIN recipe rc ON l.recipe_id = rc.id
		ORDER BY e.type, e.name
	`
	// 使用 rawEq 先接收，再手動對應 JOIN 欄位
	rows, err := r.db.QueryxContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Equipment
	for rows.Next() {
		var e model.Equipment
		var currentLot, recipeName sql.NullString
		err := rows.Scan(
			&e.ID, &e.Name, &e.Type, &e.Status, &e.CurrentLotID,
			&currentLot, &recipeName,
			&e.Utilization, &e.Temperature, &e.Pressure,
			&e.UCLTemp, &e.LCLTemp, &e.UCLPressure, &e.LCLPressure,
			&e.IsAlarm, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if currentLot.Valid {
			e.CurrentLot = &currentLot.String
		}
		if recipeName.Valid {
			e.RecipeName = &recipeName.String
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

// FindByID 依 id 取得單台設備
func (r *equipmentRepo) FindByID(ctx context.Context, id int) (*model.Equipment, error) {
	const q = `
		SELECT
			e.id, e.name, e.type, e.status, e.current_lot_id,
			l.lot_number  AS current_lot,
			rc.name       AS recipe_name,
			e.utilization, e.temperature, e.pressure,
			e.ucl_temp, e.lcl_temp, e.ucl_pressure, e.lcl_pressure,
			e.is_alarm, e.updated_at
		FROM equipment e
		LEFT JOIN lot l    ON e.current_lot_id = l.id
		LEFT JOIN recipe rc ON l.recipe_id = rc.id
		WHERE e.id = $1
	`
	row := r.db.QueryRowxContext(ctx, q, id)
	var e model.Equipment
	var currentLot, recipeName sql.NullString
	err := row.Scan(
		&e.ID, &e.Name, &e.Type, &e.Status, &e.CurrentLotID,
		&currentLot, &recipeName,
		&e.Utilization, &e.Temperature, &e.Pressure,
		&e.UCLTemp, &e.LCLTemp, &e.UCLPressure, &e.LCLPressure,
		&e.IsAlarm, &e.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if currentLot.Valid {
		e.CurrentLot = &currentLot.String
	}
	if recipeName.Valid {
		e.RecipeName = &recipeName.String
	}
	return &e, nil
}

// UpdateStatus 更新設備狀態
func (r *equipmentRepo) UpdateStatus(ctx context.Context, id int, status model.EquipmentStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE equipment SET status = $1, updated_at = $2 WHERE id = $3`,
		status, time.Now(), id,
	)
	return err
}

// FindIdleByType 找該類型第一台閒置設備（用於派工）
func (r *equipmentRepo) FindIdleByType(ctx context.Context, equipType string) (*model.Equipment, error) {
	const q = `
		SELECT id, name, type, status, current_lot_id,
		       utilization, temperature, pressure,
		       ucl_temp, lcl_temp, ucl_pressure, lcl_pressure,
		       is_alarm, updated_at
		FROM equipment
		WHERE type = $1 AND status = 'IDLE'
		ORDER BY id
		LIMIT 1
	`
	var e model.Equipment
	err := r.db.GetContext(ctx, &e, q, equipType)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &e, err
}

// AssignLot 將 Lot 指派給設備，並將設備狀態改為 RUNNING
func (r *equipmentRepo) AssignLot(ctx context.Context, equipID, lotID int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE equipment SET current_lot_id = $1, status = 'RUNNING', updated_at = $2 WHERE id = $3`,
		lotID, time.Now(), equipID,
	)
	return err
}

// ClearLot 清除設備上的 Lot，並將狀態改為 IDLE
func (r *equipmentRepo) ClearLot(ctx context.Context, equipID int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE equipment SET current_lot_id = NULL, status = 'IDLE', updated_at = $1 WHERE id = $2`,
		time.Now(), equipID,
	)
	return err
}
