// Package models models/transaction.go
package models

import "github.com/shopspring/decimal"

type TransactionType string

const (
	TransactionDeposit  TransactionType = "deposit"
	TransactionWithdraw TransactionType = "withdraw"
)

type Transaction struct {
	Base
	WalletID      uint            `gorm:"not null;index"`
	Type          TransactionType `gorm:"not null;size:20"`
	Amount        decimal.Decimal `gorm:"not null"`
	BalanceBefore decimal.Decimal `gorm:"not null"`
	BalanceAfter  decimal.Decimal `gorm:"not null"`
	Status        string          `gorm:"not null;default:'pending'"`
	TxHash        string          `gorm:"size:100;index"`
	Description   string          `gorm:"size:255"`
}
