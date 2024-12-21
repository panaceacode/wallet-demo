// Package models models/transaction.go
package models

type TransactionType string

const (
	TransactionDeposit  TransactionType = "deposit"
	TransactionWithdraw TransactionType = "withdraw"
)

type Transaction struct {
	Base
	WalletID      uint            `gorm:"not null;index"`
	Type          TransactionType `gorm:"not null;size:20"`
	Amount        float64         `gorm:"not null"`
	BalanceBefore float64         `gorm:"not null"`
	BalanceAfter  float64         `gorm:"not null"`
	Status        string          `gorm:"not null;default:'pending'"`
	TxHash        string          `gorm:"size:100;index"`
	Description   string          `gorm:"size:255"`
}
