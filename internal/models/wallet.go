package models

type Wallet struct {
	ID      uint `gorm:"primary_key;autoIncrement" json:"id"`
	Balance int  `gorm:"not null;default:0" json:"balance"`
}
