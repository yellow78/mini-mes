package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/yellow78/mini-mes/backend/internal/model"
)

type alarmRepo struct {
	db *sqlx.DB
}

// NewAlarmRepository 建立告警 repository 實例
func NewAlarmRepository(db *sqlx.DB) AlarmRepository {
	return &alarmRepo{db: db}
}

// FindAll 取得告警列表，JOIN equipment 取得設備名稱
func (r *alarmRepo) FindAll(ctx context.Context, onlyUnacked bool) ([]model.AlarmEvent, error) {
	q := `
		SELECT a.id, a.equipment_id, a.parameter, a.value, a.ucl, a.lcl,
		       a.severity, a.acknowledged, a.timestamp,
		       e.name AS equipment_name
		FROM alarm_event a
		JOIN equipment e ON a.equipment_id = e.id
	`
	if onlyUnacked {
		q += ` WHERE a.acknowledged = false`
	}
	q += ` ORDER BY a.timestamp DESC LIMIT 100`

	rows, err := r.db.QueryxContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.AlarmEvent
	for rows.Next() {
		var a model.AlarmEvent
		err := rows.Scan(
			&a.ID, &a.EquipmentID, &a.Parameter, &a.Value, &a.UCL, &a.LCL,
			&a.Severity, &a.Acknowledged, &a.Timestamp, &a.EquipmentName,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}

// Acknowledge 確認告警
func (r *alarmRepo) Acknowledge(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE alarm_event SET acknowledged = true WHERE id = $1`,
		id,
	)
	return err
}

// Create 新增告警紀錄
func (r *alarmRepo) Create(ctx context.Context, alarm *model.AlarmEvent) error {
	const q = `
		INSERT INTO alarm_event (equipment_id, parameter, value, ucl, lcl, severity)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, timestamp
	`
	return r.db.QueryRowxContext(ctx, q,
		alarm.EquipmentID, alarm.Parameter, alarm.Value,
		alarm.UCL, alarm.LCL, alarm.Severity,
	).Scan(&alarm.ID, &alarm.Timestamp)
}
