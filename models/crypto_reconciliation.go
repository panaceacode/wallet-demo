package models

import "time"

type CryptoReconciliation struct {
	Base
	WalletID       uint `gorm:"not null;index"`
	StartTime      time.Time
	EndTime        time.Time
	SystemBalance  float64
	ChainBalance   float64
	Status         ReconciliationStatus
	Difference     float64
	MismatchReason string `gorm:"type:text"`
	UnmatchedTxs   string `gorm:"type:text"` // JSON array of unmatched transaction hashes
}
