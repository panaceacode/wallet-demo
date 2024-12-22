// Package models models/wallet.go
package models

import (
	"github.com/shopspring/decimal"
)

type Wallet struct {
	Base
	UserID   uint            `gorm:"not null;index"`
	Currency string          `gorm:"not null;size:10"`
	Balance  decimal.Decimal `gorm:"not null;default:0"`
}
