// Package controllers controllers/wallet_controller.go
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/panaceacode/wallet-demo/services"
	"net/http"
	"time"
)

type WalletController struct {
	walletService         *services.WalletService
	reconciliationService *services.ReconciliationService
}

func NewWalletController(walletService *services.WalletService, reconciliationService *services.ReconciliationService) *WalletController {
	return &WalletController{
		walletService: walletService,
	}
}

type CreateWalletRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	TxHash string  `json:"tx_hash" binding:"required"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	TxHash string  `json:"tx_hash" binding:"required"`
}

func (c *WalletController) CreateWallet(ctx *gin.Context) {
	var req CreateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := c.walletService.CreateWallet(req.UserID, req.Currency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, wallet)
}

func (c *WalletController) Deposit(ctx *gin.Context) {
	walletID := ctx.GetUint("wallet_id")
	var req DepositRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.walletService.Deposit(walletID, req.Amount, req.TxHash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deposit successful"})
}

func (c *WalletController) Withdraw(ctx *gin.Context) {
	walletID := ctx.GetUint("wallet_id")
	var req WithdrawRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.walletService.Withdraw(walletID, req.Amount, req.TxHash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "withdrawal successful"})
}

func (c *WalletController) GetTransactions(ctx *gin.Context) {
	walletID := ctx.GetUint("wallet_id")
	page := ctx.GetInt("page")
	pageSize := ctx.GetInt("page_size")

	transactions, err := c.walletService.GetTransactions(walletID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

type PerformReconciliationRequest struct {
	StartTime       time.Time `json:"start_time" binding:"required"`
	EndTime         time.Time `json:"end_time" binding:"required"`
	ExternalBalance float64   `json:"external_balance" binding:"required"`
}

func (c *WalletController) PerformReconciliation(ctx *gin.Context) {
	walletID := ctx.GetUint("wallet_id")
	var req PerformReconciliationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reconciliation, err := c.reconciliationService.PerformReconciliation(
		walletID,
		req.StartTime,
		req.EndTime,
		req.ExternalBalance,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reconciliation)
}

func (c *WalletController) GetReconciliationHistory(ctx *gin.Context) {
	walletID := ctx.GetUint("wallet_id")
	page := ctx.GetInt("page")
	pageSize := ctx.GetInt("page_size")

	reconciliations, err := c.reconciliationService.GetReconciliationHistory(walletID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reconciliations)
}

func (c *WalletController) GetReconciliationDetail(ctx *gin.Context) {
	reconciliationID := ctx.GetUint("reconciliation_id")

	reconciliation, transactions, err := c.reconciliationService.GetReconciliationDetail(reconciliationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reconciliation": reconciliation,
		"transactions":   transactions,
	})
}
