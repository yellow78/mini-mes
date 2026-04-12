package simulator

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/yellow78/mini-mes/backend/internal/model"
	"github.com/yellow78/mini-mes/backend/internal/service"
)

// Broadcaster WebSocket 廣播介面
type Broadcaster interface {
	Broadcast(event string, payload any)
}

// Simulator 產線模擬器（Demo 用，定時廣播即時事件）
type Simulator struct {
	equipSvc *service.EquipmentService
	hub      Broadcaster
	rng      *rand.Rand
}

func NewSimulator(equipSvc *service.EquipmentService, hub Broadcaster) *Simulator {
	return &Simulator{
		equipSvc: equipSvc,
		hub:      hub,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Start 啟動模擬器（背景 goroutine）
func (s *Simulator) Start(ctx context.Context) {
	go s.runSPCTicker(ctx)
	go s.runStatusTicker(ctx)
	log.Println("[Simulator] 產線模擬器啟動（SPC: 每6秒 / 狀態: 每25秒）")
}

// runSPCTicker 每 6 秒模擬一次 SPC 量測，10% 機率超標觸發告警
func (s *Simulator) runSPCTicker(ctx context.Context) {
	ticker := time.NewTicker(6 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.simulateSPC(ctx)
		}
	}
}

// runStatusTicker 每 25 秒隨機切換一台設備的 RUNNING/IDLE 狀態
func (s *Simulator) runStatusTicker(ctx context.Context) {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.simulateStatusChange(ctx)
		}
	}
}

// simulateSPC 模擬量測，超過 UCL/LCL 即廣播 spc_alarm
func (s *Simulator) simulateSPC(ctx context.Context) {
	equipments, err := s.equipSvc.GetAll(ctx)
	if err != nil || len(equipments) == 0 {
		return
	}

	// 只對 RUNNING 設備進行量測
	var running []model.Equipment
	for _, e := range equipments {
		if e.Status == model.StatusRunning {
			running = append(running, e)
		}
	}
	if len(running) == 0 {
		return
	}

	eq := running[s.rng.Intn(len(running))]

	// 隨機選擇量測參數（溫度或壓力）
	var parameter string
	var value, ucl, lcl float64

	if s.rng.Intn(2) == 0 {
		parameter = "temperature"
		ucl = eq.UCLTemp
		lcl = eq.LCLTemp
	} else {
		parameter = "pressure"
		ucl = eq.UCLPressure
		lcl = eq.LCLPressure
	}

	mid := (ucl + lcl) / 2
	// 10% 機率生成超標值（UCL 以上 2–6%），其餘在正常範圍波動
	if s.rng.Float64() < 0.10 {
		value = ucl * (1.02 + s.rng.Float64()*0.04)
	} else {
		value = mid + (s.rng.Float64()-0.5)*0.1*(ucl-lcl)
	}

	if value > ucl || value < lcl {
		s.hub.Broadcast("spc_alarm", map[string]any{
			"equipment_id":   eq.ID,
			"equipment_name": eq.Name,
			"parameter":      parameter,
			"value":          value,
			"ucl":            ucl,
			"lcl":            lcl,
		})
		log.Printf("[Simulator] SPC 告警：%s %s=%.2f (UCL=%.2f LCL=%.2f)",
			eq.Name, parameter, value, ucl, lcl)
	}
}

// simulateStatusChange 隨機切換一台設備的 RUNNING/IDLE 狀態並廣播
func (s *Simulator) simulateStatusChange(ctx context.Context) {
	equipments, err := s.equipSvc.GetAll(ctx)
	if err != nil || len(equipments) == 0 {
		return
	}

	// 只切換 RUNNING 或 IDLE 設備
	var candidates []model.Equipment
	for _, e := range equipments {
		if e.Status == model.StatusRunning || e.Status == model.StatusIdle {
			candidates = append(candidates, e)
		}
	}
	if len(candidates) == 0 {
		return
	}

	eq := candidates[s.rng.Intn(len(candidates))]
	var newStatus model.EquipmentStatus
	if eq.Status == model.StatusRunning {
		newStatus = model.StatusIdle
	} else {
		newStatus = model.StatusRunning
	}

	if err := s.equipSvc.UpdateStatus(ctx, eq.ID, newStatus); err != nil {
		log.Printf("[Simulator] 狀態切換失敗（%s）: %v", eq.Name, err)
		return
	}

	s.hub.Broadcast("equipment_status_changed", map[string]any{
		"equipment_id": eq.ID,
		"status":       newStatus,
	})
	log.Printf("[Simulator] 狀態切換：%s %s → %s", eq.Name, eq.Status, newStatus)
}
