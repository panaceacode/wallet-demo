package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/panaceacode/wallet-demo/db"
	"github.com/panaceacode/wallet-demo/services"
	"net/http"
)

// 查询钱包余额
func GetBalanceHandler(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "1")
	currency := c.DefaultQuery("currency", "BTC")

	var wallet db.Wallet
	err := db.DB.QueryRow(`
		SELECT user_id, currency, balance, created_at, updated_at
		FROM wallet WHERE user_id = ? AND currency = ?
	`, userID, currency).Scan(&wallet.UserID, &wallet.Currency, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

// 充值入账接口
func DepositHandler(c *gin.Context) {
	var request struct {
		UserID   int     `json:"user_id"`
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 调用 deposit 函数
	err := services.Deposit(request.UserID, request.Currency, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// 提币接口
func WithdrawHandler(c *gin.Context) {
	var req struct {
		UserID   int     `json:"user_id"`
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := services.Withdraw(req.UserID, req.Currency, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdraw successful"})
}
