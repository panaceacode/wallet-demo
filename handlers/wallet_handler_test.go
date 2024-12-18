package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/panaceacode/wallet-demo/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup() {
	db.InitDB() // 初始化内存数据库
}

func teardown() {
	db.CloseDB() // 关闭数据库
}

func TestGetBalanceHandler(t *testing.T) {
	setup()
	defer teardown()

	// Insert a wallet balance for testing
	_, err := db.DB.Exec("INSERT INTO wallet (user_id, currency, balance) VALUES (?, ?, ?)", 1, "BTC", 100)
	if err != nil {
		t.Fatalf("Failed to setup test data: %v", err)
	}

	// Create a request to fetch balance
	req, err := http.NewRequest("GET", "/balance?user_id=1&currency=BTC", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// Setup Gin engine and register the route
	router := gin.Default()
	router.GET("/balance", GetBalanceHandler)

	// Serve the HTTP request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	// Check the response body
	var wallet db.Wallet
	err = json.NewDecoder(rr.Body).Decode(&wallet)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if wallet.Balance != 100 {
		t.Fatalf("Expected balance 100, got %f", wallet.Balance)
	}
}

func TestDepositHandler(t *testing.T) {
	setup()
	defer teardown()

	// Create a deposit request payload
	payload := map[string]interface{}{
		"user_id":  1,
		"currency": "BTC",
		"amount":   100,
	}
	jsonPayload, _ := json.Marshal(payload)

	// Create a request to deposit
	req, err := http.NewRequest("POST", "/deposit", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// Setup Gin engine and register the route
	router := gin.Default()
	router.POST("/deposit", DepositHandler)

	// Serve the HTTP request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
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

func TestWithdrawHandler(t *testing.T) {
	setup()
	defer teardown()

	// First, deposit 100 BTC to the wallet
	_, err := db.DB.Exec("INSERT INTO wallet (user_id, currency, balance) VALUES (?, ?, ?)", 1, "BTC", 100)
	if err != nil {
		t.Fatalf("Failed to setup test data: %v", err)
	}

	// Create a withdraw request payload
	payload := map[string]interface{}{
		"user_id":  1,
		"currency": "BTC",
		"amount":   50,
	}
	jsonPayload, _ := json.Marshal(payload)

	// Create a request to withdraw
	req, err := http.NewRequest("POST", "/withdraw", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// Setup Gin engine and register the route
	router := gin.Default()
	router.POST("/withdraw", WithdrawHandler)

	// Serve the HTTP request
	router.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
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
