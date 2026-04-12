package simulator

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/yellow78/mini-mes/backend/internal/model"
	"github.com/yellow78/mini-mes/backend/internal/repository"
	"github.com/yellow78/mini-mes/backend/internal/service"
	"github.com/yellow78/mini-mes/backend/pkg/spcclient"
)

// Broadcaster WebSocket 廣播介面
type Broadcaster interface {
	Broadcast(event string, payload any)
}

// Simulator 產線模擬器（Demo 用）
// 定時生成 SPC 量測值，透過 Python 引擎分析後寫入 DB 並廣播告警
type Simulator struct {
	equipSvc  *service.EquipmentService
	spcRepo   repository.SpcRepository
	alarmRepo repository.AlarmRepository
	spcClient *spcclient.Client
	hub       Broadcaster
	rng       *rand.Rand
}

func NewSimulator(
	equipSvc *service.EquipmentService,
	spcRepo repository.SpcRepository,
	alarmRepo repository.AlarmRepository,
	spcClient *spcclient.Client,
	hub Broadcaster,
) *Simulator {
	return &Simulator{
		equipSvc:  equipSvc,
		spcRepo:   spcRepo,
		alarmRepo: alarmRepo,
		spcClient: spcClient,
		hub:       hub,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Start 啟動模擬器（背景 goroutine）
func (s *Simulator) Start(ctx context.Context) {
	go s.runSPCTicker(ctx)
	go s.runStatusTicker(ctx)
	log.Println("[Simulator] 產線模擬器啟動（SPC: 每8秒 / 狀態: 每25秒）")
}

// runSPCTicker 每 8 秒模擬一次 SPC 量測
func (s *Simulator) runSPCTicker(ctx context.Context) {
	ticker := time.NewTicker(8 * time.Second)
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

// simulateSPC 生成新量測值 → Python 分析 → 寫 DB → 必要時廣播告警
func (s *Simulator) simulateSPC(ctx context.Context) {
	equipments, err := s.equipSvc.GetAll(ctx)
	if err != nil || len(equipments) == 0 {
		return
	}

	// 只對 RUNNING 設備量測
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

	// 選擇量測參數（溫度或壓力）
	var parameter string
	var newValue, ucl, lcl float64
	if s.rng.Intn(2) == 0 {
		parameter = "temperature"
		ucl, lcl = eq.UCLTemp, eq.LCLTemp
	} else {
		parameter = "pressure"
		ucl, lcl = eq.UCLPressure, eq.LCLPressure
	}

	// 生成新量測值（90% 正常範圍，10% 超標）
	mid := (ucl + lcl) / 2
	if s.rng.Float64() < 0.10 {
		newValue = ucl * (1.02 + s.rng.Float64()*0.04)
	} else {
		newValue = mid + (s.rng.Float64()-0.5)*0.1*(ucl-lcl)
	}

	// 取得近 20 筆歷史紀錄，組成分析序列
	history, err := s.spcRepo.FindByEquipment(ctx, eq.ID, 20)
	if err != nil {
		log.Printf("[Simulator] 讀取 SPC 歷史失敗: %v", err)
		return
	}

	// 歷史紀錄由新到舊，反轉後加上新值
	values := make([]float64, 0, len(history)+1)
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Parameter == parameter {
			values = append(values, history[i].Value)
		}
	}
	values = append(values, newValue)

	// 呼叫 Python SPC 引擎分析
	result, err := s.spcClient.Analyze(ctx, spcclient.AnalyzeRequest{
		EquipmentID: eq.ID,
		Parameter:   parameter,
		Values:      values,
		UCL:         ucl,
		LCL:         lcl,
	})
	if err != nil {
		// Python 不可用時降級：只做簡單閾值判斷
		log.Printf("[Simulator] Python SPC 不可用，降級處理: %v", err)
		s.fallbackSPC(ctx, eq, parameter, newValue, ucl, lcl)
		return
	}

	// 寫入 SPC 紀錄
	record := &model.SpcRecord{
		EquipmentID: eq.ID,
		Parameter:   parameter,
		Value:       newValue,
		UCL:         ucl,
		LCL:         lcl,
		IsAlarm:     result.IsAlarm,
	}
	if err := s.spcRepo.Create(ctx, record); err != nil {
		log.Printf("[Simulator] 寫入 SPC 紀錄失敗: %v", err)
	}

	if !result.IsAlarm {
		return
	}

	// 寫入告警事件
	alarmEvent := &model.AlarmEvent{
		EquipmentID: eq.ID,
		Parameter:   parameter,
		Value:       newValue,
		UCL:         ucl,
		LCL:         lcl,
		Severity:    result.Severity,
	}
	if err := s.alarmRepo.Create(ctx, alarmEvent); err != nil {
		log.Printf("[Simulator] 寫入告警失敗: %v", err)
	}

	// 廣播 spc_alarm 事件
	s.hub.Broadcast("spc_alarm", map[string]any{
		"equipment_id":   eq.ID,
		"equipment_name": eq.Name,
		"parameter":      parameter,
		"value":          newValue,
		"ucl":            ucl,
		"lcl":            lcl,
	})
	log.Printf("[Simulator] SPC 告警（%s）：%s %s=%.2f 違規：%v",
		result.Severity, eq.Name, parameter, newValue, result.Violations)
}

// fallbackSPC Python 不可用時的簡單閾值判斷
func (s *Simulator) fallbackSPC(ctx context.Context, eq model.Equipment,
	parameter string, value, ucl, lcl float64,
) {
	isAlarm := value > ucl || value < lcl
	record := &model.SpcRecord{
		EquipmentID: eq.ID,
		Parameter:   parameter,
		Value:       value,
		UCL:         ucl,
		LCL:         lcl,
		IsAlarm:     isAlarm,
	}
	if err := s.spcRepo.Create(ctx, record); err != nil {
		log.Printf("[Simulator] fallback SPC 寫入失敗: %v", err)
	}

	if !isAlarm {
		return
	}

	severity := "WARNING"
	if value > ucl*1.05 || value < lcl*0.95 {
		severity = "CRITICAL"
	}

	alarmEvent := &model.AlarmEvent{
		EquipmentID: eq.ID,
		Parameter:   parameter,
		Value:       value,
		UCL:         ucl,
		LCL:         lcl,
		Severity:    severity,
	}
	if err := s.alarmRepo.Create(ctx, alarmEvent); err != nil {
		log.Printf("[Simulator] fallback 告警寫入失敗: %v", err)
	}

	s.hub.Broadcast("spc_alarm", map[string]any{
		"equipment_id":   eq.ID,
		"equipment_name": eq.Name,
		"parameter":      parameter,
		"value":          value,
		"ucl":            ucl,
		"lcl":            lcl,
	})
}

// simulateStatusChange 隨機切換一台設備的 RUNNING/IDLE 狀態並廣播
func (s *Simulator) simulateStatusChange(ctx context.Context) {
	equipments, err := s.equipSvc.GetAll(ctx)
	if err != nil || len(equipments) == 0 {
		return
	}

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
