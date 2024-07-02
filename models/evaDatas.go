package models

import "time"

// EvaData represents the structure of the newborn_datas table
type EvaData struct {
	ID        uint       `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TrialCode string     `json:"trial_code"`
	Thermal   float64    `json:"thermal"`
	Notes     string     `json:"notes"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// TableName returns the name of the newborn_datas table
func (*EvaData) TableName() string {
	return "newborn_datas"
}
