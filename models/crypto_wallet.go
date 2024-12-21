package models

type Network string

const (
	NetworkBTC  Network = "BTC"
	NetworkETH  Network = "ETH"
	NetworkBSC  Network = "BSC"
	NetworkTRON Network = "TRON"
)

type CryptoWallet struct {
	Base
	UserID      uint    `gorm:"not null;index"`
	Network     Network `gorm:"size:10;not null"`
	Address     string  `gorm:"size:100;not null;uniqueIndex"`
	Balance     float64 `gorm:"not null;default:0"`
	PrivateKey  string  `gorm:"size:255"`  // Consider encryption
	AddressPath string  `gorm:"size:50"`   // BIP44 derivation path
	ExtraData   string  `gorm:"type:text"` // Network-specific data
}
