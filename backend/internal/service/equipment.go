package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/yellow78/mini-mes/backend/internal/model"
	"github.com/yellow78/mini-mes/backend/internal/repository"
)

// EquipmentService 設備業務邏輯
type EquipmentService struct {
	repo repository.EquipmentRepository
}

func NewEquipmentService(repo repository.EquipmentRepository) *EquipmentService {
	return &EquipmentService{repo: repo}
}

// GetAll 取得所有設備
func (s *EquipmentService) GetAll(ctx context.Context) ([]model.Equipment, error) {
	return s.repo.FindAll(ctx)
}

// GetByID 取得單台設備，找不到回傳 nil, nil
func (s *EquipmentService) GetByID(ctx context.Context, id int) (*model.Equipment, error) {
	return s.repo.FindByID(ctx, id)
}

// UpdateStatus 更新設備狀態（含狀態轉換驗證）
func (s *EquipmentService) UpdateStatus(ctx context.Context, id int, newStatus model.EquipmentStatus) error {
	eq, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if eq == nil {
		return fmt.Errorf("equipment %d not found", id)
	}
	if !model.IsValidTransition(eq.Status, newStatus) {
		return fmt.Errorf("invalid status transition: %s → %s", eq.Status, newStatus)
	}
	return s.repo.UpdateStatus(ctx, id, newStatus)
}

// Hold 設備：強制將狀態改為 DOWN，不受轉換規則限制
func (s *EquipmentService) Hold(ctx context.Context, id int) error {
	eq, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if eq == nil {
		return fmt.Errorf("equipment %d not found", id)
	}
	// Hold 操作清除 Lot，設備進入 DOWN
	if err := s.repo.ClearLot(ctx, id); err != nil {
		return err
	}
	return s.repo.UpdateStatus(ctx, id, model.StatusDown)
}

// GetGroups 取得設備群組（依 type 分群，alarm 設備置頂）
func (s *EquipmentService) GetGroups(ctx context.Context) ([]model.EquipmentGroup, error) {
	equipments, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 依類型分群
	groupMap := make(map[string]*model.EquipmentGroup)
	order := []string{"CVD", "Etch", "CMP", "Diffusion"}

	for i := range equipments {
		eq := equipments[i]
		g, ok := groupMap[eq.Type]
		if !ok {
			g = &model.EquipmentGroup{Type: eq.Type}
			groupMap[eq.Type] = g
		}
		g.Equipments = append(g.Equipments, eq)
		if eq.IsAlarm {
			g.AlarmCount++
		}
		switch eq.Status {
		case model.StatusRunning:
			g.StatusCount.Running++
		case model.StatusIdle:
			g.StatusCount.Idle++
		case model.StatusDown:
			g.StatusCount.Down++
		case model.StatusPM:
			g.StatusCount.PM++
		}
	}

	// alarm 設備置頂，計算群組稼動率
	var groups []model.EquipmentGroup
	for _, t := range order {
		g, ok := groupMap[t]
		if !ok {
			continue
		}
		sort.Slice(g.Equipments, func(i, j int) bool {
			// alarm 排前面；同為 alarm 或同非 alarm 則保持原序
			return g.Equipments[i].IsAlarm && !g.Equipments[j].IsAlarm
		})
		// 計算群組平均稼動率（只算 RUNNING 設備）
		running := 0
		total := 0.0
		for _, e := range g.Equipments {
			if e.Status == model.StatusRunning {
				running++
				total += e.Utilization
			}
		}
		if running > 0 {
			g.Utilization = total / float64(running)
		}
		groups = append(groups, *g)
	}
	return groups, nil
}
