package service

import (
	"context"
	"fmt"

	"github.com/yellow78/mini-mes/backend/internal/model"
	"github.com/yellow78/mini-mes/backend/internal/repository"
)

// DispatchService 派工業務邏輯
type DispatchService struct {
	equipRepo repository.EquipmentRepository
	lotRepo   repository.LotRepository
}

func NewDispatchService(
	equipRepo repository.EquipmentRepository,
	lotRepo repository.LotRepository,
) *DispatchService {
	return &DispatchService{equipRepo: equipRepo, lotRepo: lotRepo}
}

// DispatchResult 派工結果
type DispatchResult struct {
	LotID       int `json:"lot_id"`
	EquipmentID int `json:"equipment_id"`
}

// Dispatch 將 Lot 指派給可用設備
// 邏輯：依 Lot 的 Recipe 找對應類型的第一台 IDLE 設備
func (s *DispatchService) Dispatch(ctx context.Context, lotID int) (*DispatchResult, error) {
	lot, err := s.lotRepo.FindByID(ctx, lotID)
	if err != nil {
		return nil, err
	}
	if lot == nil {
		return nil, fmt.Errorf("lot %d not found", lotID)
	}
	if lot.Status != model.LotQueued {
		return nil, fmt.Errorf("lot %d is not in QUEUED status (current: %s)", lotID, lot.Status)
	}

	// 從 recipe_name 推斷設備類型（需要 JOIN 查到 recipe，這裡簡化為依 recipe_id 對應）
	// Phase 2 以 recipe 的 equipment_type 欄位決定，目前透過 lot.RecipeName 前綴判斷
	equipType, err := s.resolveEquipmentType(ctx, lot)
	if err != nil {
		return nil, err
	}

	// 找第一台閒置設備
	equip, err := s.equipRepo.FindIdleByType(ctx, equipType)
	if err != nil {
		return nil, err
	}
	if equip == nil {
		return nil, fmt.Errorf("no idle %s equipment available", equipType)
	}

	// 指派 Lot 給設備
	if err := s.equipRepo.AssignLot(ctx, equip.ID, lot.ID); err != nil {
		return nil, err
	}
	// 更新 Lot 狀態為 RUNNING
	if err := s.lotRepo.UpdateStatus(ctx, lot.ID, model.LotRunning); err != nil {
		return nil, err
	}

	return &DispatchResult{LotID: lot.ID, EquipmentID: equip.ID}, nil
}

// resolveEquipmentType 依 recipe 名稱前綴決定設備類型
func (s *DispatchService) resolveEquipmentType(_ context.Context, lot *model.Lot) (string, error) {
	if lot.RecipeName == nil {
		return "", fmt.Errorf("lot %d has no recipe", lot.ID)
	}
	name := *lot.RecipeName
	switch {
	case len(name) >= 3 && name[:3] == "CVD":
		return "CVD", nil
	case len(name) >= 4 && name[:4] == "Etch":
		return "Etch", nil
	case len(name) >= 3 && name[:3] == "CMP":
		return "CMP", nil
	case len(name) >= 4 && name[:4] == "Diff":
		return "Diffusion", nil
	default:
		return "", fmt.Errorf("cannot determine equipment type from recipe: %s", name)
	}
}
