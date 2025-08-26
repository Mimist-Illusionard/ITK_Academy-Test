package models

type Wallet struct {
	ID      uint `gorm:"primaryKey"`
	Balance int
}
