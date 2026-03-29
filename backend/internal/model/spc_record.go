package model

import "time"

type SpcRecord struct {
	ID          int       `db:"id"           json:"id"`
	EquipmentID int       `db:"equipment_id" json:"equipment_id"`
	Parameter   string    `db:"parameter"    json:"parameter"`
	Value       float64   `db:"value"        json:"value"`
	UCL         float64   `db:"ucl"          json:"ucl"`
	LCL         float64   `db:"lcl"          json:"lcl"`
	IsAlarm     bool      `db:"is_alarm"     json:"is_alarm"`
	Timestamp   time.Time `db:"timestamp"    json:"timestamp"`
}

type AlarmEvent struct {
	ID          int       `db:"id"           json:"id"`
	EquipmentID int       `db:"equipment_id" json:"equipment_id"`
	Parameter   string    `db:"parameter"    json:"parameter"`
	Value       float64   `db:"value"        json:"value"`
	UCL         float64   `db:"ucl"          json:"ucl"`
	LCL         float64   `db:"lcl"          json:"lcl"`
	Severity    string    `db:"severity"     json:"severity"`
	Acknowledged bool     `db:"acknowledged" json:"acknowledged"`
	Timestamp   time.Time `db:"timestamp"    json:"timestamp"`
	// JOIN 欄位
	EquipmentName string `db:"-" json:"equipment_name"`
}
