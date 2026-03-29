package model

import "time"

type Recipe struct {
	ID             int       `db:"id"              json:"id"`
	Name           string    `db:"name"            json:"name"`
	EquipmentType  string    `db:"equipment_type"  json:"equipment_type"`
	TargetTemp     float64   `db:"target_temp"     json:"target_temp"`
	TargetPressure float64   `db:"target_pressure" json:"target_pressure"`
	DurationMin    int       `db:"duration_min"    json:"duration_min"`
	CreatedAt      time.Time `db:"created_at"      json:"created_at"`
}
