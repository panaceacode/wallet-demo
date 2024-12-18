package services

import (
	"fmt"
	"github.com/panaceacode/wallet-demo/db"
	"time"
)

// 充值入账
func Deposit(userID int, currency string, amount float64) error {
	// 检查钱包是否存在
	var currentBalance float64
	err := db.DB.QueryRow("SELECT balance FROM wallet WHERE user_id = ? AND currency = ?", userID, currency).Scan(&currentBalance)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return fmt.Errorf("failed to check balance: %v", err)
	}

	if err == nil { // 如果钱包已经存在，则更新余额
		_, err = db.DB.Exec("UPDATE wallet SET balance = balance + ?, updated_at = CURRENT_TIMESTAMP WHERE user_id = ? AND currency = ?", amount, userID, currency)
	} else { // 钱包不存在，插入新记录
		_, err = db.DB.Exec("INSERT INTO wallet (user_id, currency, balance, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", userID, currency, amount)
	}

	if err != nil {
		return fmt.Errorf("failed to deposit: %v", err)
	}

	// 记录交易
	_, err = db.DB.Exec("INSERT INTO transaction_record (user_id, currency, amount, status, created_at, updated_at) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)", userID, currency, amount, "success")
	if err != nil {
		return fmt.Errorf("failed to record transaction: %v", err)
	}

	return nil
}

// 提币
func Withdraw(userID int, currency string, amount float64) error {
	// 检查余额是否足够
	var balance float64
	err := db.DB.QueryRow(`
		SELECT balance FROM wallet WHERE user_id = ? AND currency = ?
	`, userID, currency).Scan(&balance)

	if err != nil {
		return fmt.Errorf("failed to query balance: %w", err)
	}
	if balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	now := time.Now()
	// 更新钱包余额，并设置更新时间
	_, err = db.DB.Exec(`
		UPDATE wallet SET balance = balance - ?, updated_at = ?
		WHERE user_id = ? AND currency = ?
	`, amount, now, userID, currency)

	return err
}
