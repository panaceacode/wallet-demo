package models

type CryptoTransaction struct {
	Base
	WalletID      uint            `gorm:"not null;index"`
	Type          TransactionType `gorm:"not null;size:20"`
	Network       Network         `gorm:"size:10;not null"`
	FromAddress   string          `gorm:"size:100;not null"`
	ToAddress     string          `gorm:"size:100;not null"`
	Amount        float64         `gorm:"not null"`
	Status        string          `gorm:"not null;default:'pending'"`
	TxHash        string          `gorm:"size:100;index"`
	Confirmations int             `gorm:"default:0"`
	BlockNumber   uint64          `gorm:"default:0"`
	GasPrice      string          `gorm:"size:50"`
	GasUsed       uint64          `gorm:"default:0"`
	Raw           string          `gorm:"type:text"`
}
