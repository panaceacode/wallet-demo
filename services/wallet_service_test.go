package services

import (
	"github.com/panaceacode/wallet-demo/db"
	"testing"
)

func setup() {
	db.InitDB() // 初始化内存数据库
}

func teardown() {
	db.CloseDB() // 关闭数据库
}

func TestDeposit(t *testing.T) {
	setup()
	defer teardown()

	// Test deposit 100 BTC
	err := Deposit(1, "BTC", 100)
	if err != nil {
		t.Fatalf("Deposit failed: %v", err)
	}

	// Verify balance after deposit
	var balance float64
	err = db.DB.QueryRow("SELECT balance FROM wallet WHERE user_id = ? AND currency = ?", 1, "BTC").Scan(&balance)
	if err != nil {
		t.Fatalf("Failed to fetch balance: %v", err)
	}

	if balance != 100 {
		t.Fatalf("Expected balance 100, got %f", balance)
	}
}

func TestWithdraw(t *testing.T) {
	setup()
	defer teardown()

	// First deposit 100 BTC
	err := Deposit(1, "BTC", 100)
	if err != nil {
		t.Fatalf("Deposit failed: %v", err)
	}

	// Test withdraw 50 BTC
	err = Withdraw(1, "BTC", 50)
	if err != nil {
		t.Fatalf("Withdraw failed: %v", err)
	}

	// Verify balance after withdraw
	var balance float64
	err = db.DB.QueryRow("SELECT balance FROM wallet WHERE user_id = ? AND currency = ?", 1, "BTC").Scan(&balance)
	if err != nil {
		t.Fatalf("Failed to fetch balance: %v", err)
	}

	if balance != 50 {
		t.Fatalf("Expected balance 50, got %f", balance)
	}
}

func TestWithdrawInsufficientBalance(t *testing.T) {
	setup()
	defer teardown()

	// Test withdraw with insufficient balance
	err := Withdraw(1, "BTC", 100) // Should fail as no balance
	if err == nil {
		t.Fatalf("Expected error for insufficient balance, but got none")
	}
}

func TestDepositAndTransaction(t *testing.T) {
	setup()
	defer teardown()

	// Test deposit and check transaction table
	err := Deposit(1, "BTC", 200)
	if err != nil {
		t.Fatalf("Deposit failed: %v", err)
	}

	// Verify transaction record
	var txCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM transaction_record WHERE user_id = ? AND currency = ? AND status = ?", 1, "BTC", "success").Scan(&txCount)
	if err != nil {
		t.Fatalf("Failed to query transaction: %v", err)
	}

	if txCount != 1 {
		t.Fatalf("Expected 1 transaction, got %d", txCount)
	}
}
