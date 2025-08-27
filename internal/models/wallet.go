package models

import "github.com/google/uuid"

type Wallet struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	Balance int       `gorm:"not null;default:0" json:"balance"`
}
