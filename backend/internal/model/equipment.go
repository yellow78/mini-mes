package model

import "time"

type EquipmentStatus string

const (
	StatusIdle    EquipmentStatus = "IDLE"
	StatusRunning EquipmentStatus = "RUNNING"
	StatusDown    EquipmentStatus = "DOWN"
	StatusPM      EquipmentStatus = "PM"
)

// 合法的狀態轉換規則
var ValidTransitions = map[EquipmentStatus][]EquipmentStatus{
	StatusIdle:    {StatusRunning, StatusPM, StatusDown},
	StatusRunning: {StatusIdle, StatusDown},
	StatusDown:    {StatusPM},
	StatusPM:      {StatusIdle},
}

// IsValidTransition 驗證狀態轉換是否合法
func IsValidTransition(from, to EquipmentStatus) bool {
	targets, ok := ValidTransitions[from]
	if !ok {
		return false
	}
	for _, t := range targets {
		if t == to {
			return true
		}
	}
	return false
}

type Equipment struct {
	ID           int             `db:"id"             json:"id"`
	Name         string          `db:"name"           json:"name"`
	Type         string          `db:"type"           json:"type"`
	Status       EquipmentStatus `db:"status"         json:"status"`
	CurrentLotID *int            `db:"current_lot_id" json:"current_lot_id"`
	CurrentLot   *string         `db:"-"              json:"current_lot"`   // JOIN 欄位
	RecipeName   *string         `db:"-"              json:"recipe_name"`
	Utilization  float64         `db:"utilization"    json:"utilization"`
	Temperature  float64         `db:"temperature"    json:"temperature"`
	Pressure     float64         `db:"pressure"       json:"pressure"`
	UCLTemp      float64         `db:"ucl_temp"       json:"ucl_temp"`
	LCLTemp      float64         `db:"lcl_temp"       json:"lcl_temp"`
	UCLPressure  float64         `db:"ucl_pressure"   json:"ucl_pressure"`
	LCLPressure  float64         `db:"lcl_pressure"   json:"lcl_pressure"`
	IsAlarm      bool            `db:"is_alarm"       json:"is_alarm"`
	UpdatedAt    time.Time       `db:"updated_at"     json:"updated_at"`
}

// EquipmentGroup 群組折疊用
type EquipmentGroup struct {
	Type        string      `json:"type"`
	Equipments  []Equipment `json:"equipments"`
	AlarmCount  int         `json:"alarm_count"`
	Utilization float64     `json:"utilization"`
	StatusCount struct {
		Running int `json:"running"`
		Idle    int `json:"idle"`
		Down    int `json:"down"`
		PM      int `json:"pm"`
	} `json:"status_count"`
}
