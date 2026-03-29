package service

import (
	"context"
	"fmt"

	"github.com/yellow78/mini-mes/backend/internal/model"
	"github.com/yellow78/mini-mes/backend/internal/repository"
)

// LotService Lot 業務邏輯
type LotService struct {
	repo repository.LotRepository
}

func NewLotService(repo repository.LotRepository) *LotService {
	return &LotService{repo: repo}
}

// GetAll 取得所有 Lot
func (s *LotService) GetAll(ctx context.Context) ([]model.Lot, error) {
	return s.repo.FindAll(ctx)
}

// GetByID 取得單筆 Lot
func (s *LotService) GetByID(ctx context.Context, id int) (*model.Lot, error) {
	return s.repo.FindByID(ctx, id)
}

// Create 建立新 Lot，預設狀態 QUEUED
func (s *LotService) Create(ctx context.Context, lot *model.Lot) error {
	if lot.LotNumber == "" {
		return fmt.Errorf("lot_number is required")
	}
	if lot.WaferCount <= 0 {
		lot.WaferCount = 25
	}
	lot.Status = model.LotQueued
	return s.repo.Create(ctx, lot)
}

// UpdateStatus 更新 Lot 狀態
func (s *LotService) UpdateStatus(ctx context.Context, id int, status model.LotStatus) error {
	lot, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if lot == nil {
		return fmt.Errorf("lot %d not found", id)
	}
	return s.repo.UpdateStatus(ctx, id, status)
}
