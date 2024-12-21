// Package models models/wallet.go
package models

type Wallet struct {
	Base
	UserID   uint    `gorm:"not null;index"`
	Currency string  `gorm:"not null;size:10"`
	Balance  float64 `gorm:"not null;default:0"`
}
