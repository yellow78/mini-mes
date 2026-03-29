package repository

import (
	"context"

	"github.com/yellow78/mini-mes/backend/internal/model"
)

// EquipmentRepository 設備 DB 存取介面（供 service mock 使用）
type EquipmentRepository interface {
	FindAll(ctx context.Context) ([]model.Equipment, error)
	FindByID(ctx context.Context, id int) (*model.Equipment, error)
	UpdateStatus(ctx context.Context, id int, status model.EquipmentStatus) error
	FindIdleByType(ctx context.Context, equipType string) (*model.Equipment, error)
	AssignLot(ctx context.Context, equipID, lotID int) error
	ClearLot(ctx context.Context, equipID int) error
}

// LotRepository Lot DB 存取介面
type LotRepository interface {
	FindAll(ctx context.Context) ([]model.Lot, error)
	FindByID(ctx context.Context, id int) (*model.Lot, error)
	Create(ctx context.Context, lot *model.Lot) error
	UpdateStatus(ctx context.Context, id int, status model.LotStatus) error
}

// AlarmRepository 告警 DB 存取介面
type AlarmRepository interface {
	FindAll(ctx context.Context, onlyUnacked bool) ([]model.AlarmEvent, error)
	Acknowledge(ctx context.Context, id int) error
	Create(ctx context.Context, alarm *model.AlarmEvent) error
}

// SpcRepository SPC 紀錄 DB 存取介面
type SpcRepository interface {
	FindByEquipment(ctx context.Context, equipmentID int, limit int) ([]model.SpcRecord, error)
	Create(ctx context.Context, record *model.SpcRecord) error
}
