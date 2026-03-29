package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/yellow78/mini-mes/backend/internal/model"
)

type spcRepo struct {
	db *sqlx.DB
}

// NewSpcRepository 建立 SPC repository 實例
func NewSpcRepository(db *sqlx.DB) SpcRepository {
	return &spcRepo{db: db}
}

// FindByEquipment 取得指定設備最近 N 筆 SPC 紀錄
func (r *spcRepo) FindByEquipment(ctx context.Context, equipmentID int, limit int) ([]model.SpcRecord, error) {
	const q = `
		SELECT id, equipment_id, parameter, value, ucl, lcl, is_alarm, timestamp
		FROM spc_record
		WHERE equipment_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`
	var records []model.SpcRecord
	err := r.db.SelectContext(ctx, &records, q, equipmentID, limit)
	return records, err
}

// Create 新增 SPC 紀錄
func (r *spcRepo) Create(ctx context.Context, record *model.SpcRecord) error {
	const q = `
		INSERT INTO spc_record (equipment_id, parameter, value, ucl, lcl, is_alarm)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, timestamp
	`
	return r.db.QueryRowxContext(ctx, q,
		record.EquipmentID, record.Parameter, record.Value,
		record.UCL, record.LCL, record.IsAlarm,
	).Scan(&record.ID, &record.Timestamp)
}
