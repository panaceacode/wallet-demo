// Package models models/reconciliation.go
package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type ReconciliationStatus string

const (
	ReconciliationStatusPending  ReconciliationStatus = "pending"
	ReconciliationStatusMatched  ReconciliationStatus = "matched"
	ReconciliationStatusMismatch ReconciliationStatus = "mismatch"
)

type Reconciliation struct {
	Base
	WalletID        uint                 `gorm:"not null;index"`
	StartTime       time.Time            `gorm:"not null"`
	EndTime         time.Time            `gorm:"not null"`
	SystemBalance   decimal.Decimal      `gorm:"not null"` // 系统计算的余额
	ExternalBalance decimal.Decimal      `gorm:"not null"` // 外部系统的余额
	Status          ReconciliationStatus `gorm:"not null"`
	Difference      decimal.Decimal      `gorm:"not null"` // 差额
	Notes           string               `gorm:"type:text"`
}
