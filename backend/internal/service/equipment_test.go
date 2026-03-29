package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yellow78/mini-mes/backend/internal/model"
)

// --- Mock Repository ---

type mockEquipmentRepo struct {
	equipments []model.Equipment
	updateErr  error
	clearErr   error
}

func (m *mockEquipmentRepo) FindAll(_ context.Context) ([]model.Equipment, error) {
	return m.equipments, nil
}

func (m *mockEquipmentRepo) FindByID(_ context.Context, id int) (*model.Equipment, error) {
	for i := range m.equipments {
		if m.equipments[i].ID == id {
			return &m.equipments[i], nil
		}
	}
	return nil, nil
}

func (m *mockEquipmentRepo) UpdateStatus(_ context.Context, id int, status model.EquipmentStatus) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	for i := range m.equipments {
		if m.equipments[i].ID == id {
			m.equipments[i].Status = status
		}
	}
	return nil
}

func (m *mockEquipmentRepo) FindIdleByType(_ context.Context, t string) (*model.Equipment, error) {
	for i := range m.equipments {
		if m.equipments[i].Type == t && m.equipments[i].Status == model.StatusIdle {
			return &m.equipments[i], nil
		}
	}
	return nil, nil
}

func (m *mockEquipmentRepo) AssignLot(_ context.Context, equipID, lotID int) error {
	return nil
}

func (m *mockEquipmentRepo) ClearLot(_ context.Context, _ int) error {
	return m.clearErr
}

// --- UpdateStatus 測試 ---

func TestEquipmentService_UpdateStatus(t *testing.T) {
	t.Run("IDLE → RUNNING 合法轉換應成功", func(t *testing.T) {
		mock := &mockEquipmentRepo{
			equipments: []model.Equipment{{ID: 1, Status: model.StatusIdle}},
		}
		svc := NewEquipmentService(mock)
		err := svc.UpdateStatus(context.Background(), 1, model.StatusRunning)
		assert.NoError(t, err)
		assert.Equal(t, model.StatusRunning, mock.equipments[0].Status)
	})

	t.Run("DOWN → RUNNING 非法轉換應回傳錯誤", func(t *testing.T) {
		mock := &mockEquipmentRepo{
			equipments: []model.Equipment{{ID: 1, Status: model.StatusDown}},
		}
		svc := NewEquipmentService(mock)
		err := svc.UpdateStatus(context.Background(), 1, model.StatusRunning)
		assert.ErrorContains(t, err, "invalid status transition")
	})

	t.Run("RUNNING → IDLE 合法轉換應成功", func(t *testing.T) {
		mock := &mockEquipmentRepo{
			equipments: []model.Equipment{{ID: 2, Status: model.StatusRunning}},
		}
		svc := NewEquipmentService(mock)
		err := svc.UpdateStatus(context.Background(), 2, model.StatusIdle)
		assert.NoError(t, err)
	})

	t.Run("PM → RUNNING 非法轉換應回傳錯誤", func(t *testing.T) {
		mock := &mockEquipmentRepo{
			equipments: []model.Equipment{{ID: 3, Status: model.StatusPM}},
		}
		svc := NewEquipmentService(mock)
		err := svc.UpdateStatus(context.Background(), 3, model.StatusRunning)
		assert.ErrorContains(t, err, "invalid status transition")
	})

	t.Run("設備不存在應回傳錯誤", func(t *testing.T) {
		mock := &mockEquipmentRepo{equipments: []model.Equipment{}}
		svc := NewEquipmentService(mock)
		err := svc.UpdateStatus(context.Background(), 999, model.StatusRunning)
		assert.ErrorContains(t, err, "not found")
	})
}

// --- Hold 測試 ---

func TestEquipmentService_Hold(t *testing.T) {
	t.Run("Hold RUNNING 設備應成功", func(t *testing.T) {
		mock := &mockEquipmentRepo{
			equipments: []model.Equipment{{ID: 1, Status: model.StatusRunning}},
		}
		svc := NewEquipmentService(mock)
		err := svc.Hold(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("設備不存在應回傳錯誤", func(t *testing.T) {
		mock := &mockEquipmentRepo{equipments: []model.Equipment{}}
		svc := NewEquipmentService(mock)
		err := svc.Hold(context.Background(), 999)
		assert.ErrorContains(t, err, "not found")
	})
}

// --- GetGroups 測試 ---

func TestEquipmentService_GetGroups(t *testing.T) {
	t.Run("應依類型正確分群", func(t *testing.T) {
		mock := &mockEquipmentRepo{
			equipments: []model.Equipment{
				{ID: 1, Type: "CVD",  Status: model.StatusRunning},
				{ID: 2, Type: "CVD",  Status: model.StatusIdle},
				{ID: 3, Type: "Etch", Status: model.StatusRunning},
			},
		}
		svc := NewEquipmentService(mock)
		groups, err := svc.GetGroups(context.Background())
		require.NoError(t, err)
		assert.Len(t, groups, 2)
		assert.Equal(t, "CVD", groups[0].Type)
		assert.Equal(t, "Etch", groups[1].Type)
	})

	t.Run("Alarm 設備應排在群組最頂部", func(t *testing.T) {
		mock := &mockEquipmentRepo{
			equipments: []model.Equipment{
				{ID: 1, Type: "CVD", Name: "CVD-01", IsAlarm: false},
				{ID: 2, Type: "CVD", Name: "CVD-02", IsAlarm: true},
				{ID: 3, Type: "CVD", Name: "CVD-03", IsAlarm: false},
			},
		}
		svc := NewEquipmentService(mock)
		groups, err := svc.GetGroups(context.Background())
		require.NoError(t, err)
		require.Len(t, groups, 1)
		assert.Equal(t, "CVD-02", groups[0].Equipments[0].Name)
	})

	t.Run("AlarmCount 應正確計算", func(t *testing.T) {
		mock := &mockEquipmentRepo{
			equipments: []model.Equipment{
				{ID: 1, Type: "CVD", IsAlarm: true},
				{ID: 2, Type: "CVD", IsAlarm: true},
				{ID: 3, Type: "CVD", IsAlarm: false},
			},
		}
		svc := NewEquipmentService(mock)
		groups, err := svc.GetGroups(context.Background())
		require.NoError(t, err)
		assert.Equal(t, 2, groups[0].AlarmCount)
	})

	t.Run("StatusCount 應正確統計", func(t *testing.T) {
		mock := &mockEquipmentRepo{
			equipments: []model.Equipment{
				{ID: 1, Type: "Etch", Status: model.StatusRunning},
				{ID: 2, Type: "Etch", Status: model.StatusRunning},
				{ID: 3, Type: "Etch", Status: model.StatusIdle},
				{ID: 4, Type: "Etch", Status: model.StatusDown},
				{ID: 5, Type: "Etch", Status: model.StatusPM},
			},
		}
		svc := NewEquipmentService(mock)
		groups, err := svc.GetGroups(context.Background())
		require.NoError(t, err)
		sc := groups[0].StatusCount
		assert.Equal(t, 2, sc.Running)
		assert.Equal(t, 1, sc.Idle)
		assert.Equal(t, 1, sc.Down)
		assert.Equal(t, 1, sc.PM)
	})

	t.Run("群組順序應為 CVD Etch CMP Diffusion", func(t *testing.T) {
		mock := &mockEquipmentRepo{
			equipments: []model.Equipment{
				{ID: 1, Type: "Diffusion"},
				{ID: 2, Type: "CMP"},
				{ID: 3, Type: "Etch"},
				{ID: 4, Type: "CVD"},
			},
		}
		svc := NewEquipmentService(mock)
		groups, err := svc.GetGroups(context.Background())
		require.NoError(t, err)
		assert.Equal(t, []string{"CVD", "Etch", "CMP", "Diffusion"},
			[]string{groups[0].Type, groups[1].Type, groups[2].Type, groups[3].Type})
	})

	t.Run("設備為空時回傳空群組", func(t *testing.T) {
		mock := &mockEquipmentRepo{equipments: []model.Equipment{}}
		svc := NewEquipmentService(mock)
		groups, err := svc.GetGroups(context.Background())
		require.NoError(t, err)
		assert.Len(t, groups, 0)
	})
}

// --- IsValidTransition 測試 ---

func TestIsValidTransition(t *testing.T) {
	cases := []struct {
		from  model.EquipmentStatus
		to    model.EquipmentStatus
		valid bool
	}{
		{model.StatusIdle, model.StatusRunning, true},
		{model.StatusIdle, model.StatusPM, true},
		{model.StatusIdle, model.StatusDown, true},
		{model.StatusRunning, model.StatusIdle, true},
		{model.StatusRunning, model.StatusDown, true},
		{model.StatusRunning, model.StatusPM, false},
		{model.StatusDown, model.StatusPM, true},
		{model.StatusDown, model.StatusIdle, false},
		{model.StatusPM, model.StatusIdle, true},
		{model.StatusPM, model.StatusRunning, false},
	}
	for _, c := range cases {
		got := model.IsValidTransition(c.from, c.to)
		assert.Equal(t, c.valid, got, "%s → %s", c.from, c.to)
	}
}
