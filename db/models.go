package db

import "time"

// 钱包表
type Wallet struct {
	UserID    int       `json:"user_id"`
	Currency  string    `json:"currency"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UTXO 表
type UTXO struct {
	TxID        string    `json:"tx_id"`
	OutputIndex int       `json:"output_index"`
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	IsSpent     bool      `json:"is_spent"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// 交易记录表
type Transaction struct {
	TxID      string    `json:"tx_id"`
	UserID    int       `json:"user_id"`
	Currency  string    `json:"currency"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
