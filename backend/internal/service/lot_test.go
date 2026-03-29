package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yellow78/mini-mes/backend/internal/model"
)

// --- Mock Repository ---

type mockLotRepo struct {
	lots      []model.Lot
	createErr error
}

func (m *mockLotRepo) FindAll(_ context.Context) ([]model.Lot, error) {
	return m.lots, nil
}

func (m *mockLotRepo) FindByID(_ context.Context, id int) (*model.Lot, error) {
	for i := range m.lots {
		if m.lots[i].ID == id {
			return &m.lots[i], nil
		}
	}
	return nil, nil
}

func (m *mockLotRepo) Create(_ context.Context, lot *model.Lot) error {
	if m.createErr != nil {
		return m.createErr
	}
	lot.ID = len(m.lots) + 1
	m.lots = append(m.lots, *lot)
	return nil
}

func (m *mockLotRepo) UpdateStatus(_ context.Context, id int, status model.LotStatus) error {
	for i := range m.lots {
		if m.lots[i].ID == id {
			m.lots[i].Status = status
			return nil
		}
	}
	return nil
}

// --- Create 測試 ---

func TestLotService_Create(t *testing.T) {
	t.Run("建立 Lot 應自動設定 QUEUED 狀態", func(t *testing.T) {
		mock := &mockLotRepo{}
		svc := NewLotService(mock)
		lot := &model.Lot{LotNumber: "LOT-TEST-001", Product: "NAND", RecipeID: 1, Priority: 2, WaferCount: 25}
		err := svc.Create(context.Background(), lot)
		require.NoError(t, err)
		assert.Equal(t, model.LotQueued, lot.Status)
	})

	t.Run("lot_number 為空應回傳錯誤", func(t *testing.T) {
		mock := &mockLotRepo{}
		svc := NewLotService(mock)
		lot := &model.Lot{Product: "NAND", RecipeID: 1}
		err := svc.Create(context.Background(), lot)
		assert.ErrorContains(t, err, "lot_number is required")
	})

	t.Run("wafer_count <= 0 應自動補為 25", func(t *testing.T) {
		mock := &mockLotRepo{}
		svc := NewLotService(mock)
		lot := &model.Lot{LotNumber: "LOT-TEST-002", Product: "NAND", RecipeID: 1, WaferCount: 0}
		err := svc.Create(context.Background(), lot)
		require.NoError(t, err)
		assert.Equal(t, 25, lot.WaferCount)
	})

	t.Run("建立後 Lot 應有 ID", func(t *testing.T) {
		mock := &mockLotRepo{}
		svc := NewLotService(mock)
		lot := &model.Lot{LotNumber: "LOT-TEST-003", Product: "DRAM", RecipeID: 2, WaferCount: 25}
		err := svc.Create(context.Background(), lot)
		require.NoError(t, err)
		assert.Greater(t, lot.ID, 0)
	})
}

// --- UpdateStatus 測試 ---

func TestLotService_UpdateStatus(t *testing.T) {
	t.Run("正常更新狀態應成功", func(t *testing.T) {
		mock := &mockLotRepo{
			lots: []model.Lot{{ID: 1, Status: model.LotQueued}},
		}
		svc := NewLotService(mock)
		err := svc.UpdateStatus(context.Background(), 1, model.LotRunning)
		require.NoError(t, err)
		assert.Equal(t, model.LotRunning, mock.lots[0].Status)
	})

	t.Run("Lot 不存在應回傳錯誤", func(t *testing.T) {
		mock := &mockLotRepo{lots: []model.Lot{}}
		svc := NewLotService(mock)
		err := svc.UpdateStatus(context.Background(), 999, model.LotRunning)
		assert.ErrorContains(t, err, "not found")
	})
}

// --- GetByID 測試 ---

func TestLotService_GetByID(t *testing.T) {
	t.Run("存在的 Lot 應正確回傳", func(t *testing.T) {
		mock := &mockLotRepo{
			lots: []model.Lot{{ID: 1, LotNumber: "LOT-001"}},
		}
		svc := NewLotService(mock)
		lot, err := svc.GetByID(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, "LOT-001", lot.LotNumber)
	})

	t.Run("不存在的 Lot 回傳 nil", func(t *testing.T) {
		mock := &mockLotRepo{lots: []model.Lot{}}
		svc := NewLotService(mock)
		lot, err := svc.GetByID(context.Background(), 999)
		require.NoError(t, err)
		assert.Nil(t, lot)
	})
}
