package model

import "time"

type LotStatus string

const (
	LotQueued    LotStatus = "QUEUED"
	LotRunning   LotStatus = "RUNNING"
	LotCompleted LotStatus = "COMPLETED"
	LotOnHold    LotStatus = "ON_HOLD"
)

type Lot struct {
	ID         int       `db:"id"          json:"id"`
	LotNumber  string    `db:"lot_number"  json:"lot_number"`
	Product    string    `db:"product"     json:"product"`
	RecipeID   int       `db:"recipe_id"   json:"recipe_id"`
	RecipeName *string   `db:"-"           json:"recipe_name"`
	Priority   int       `db:"priority"    json:"priority"`
	Status     LotStatus `db:"status"      json:"status"`
	WaferCount int       `db:"wafer_count" json:"wafer_count"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
}
