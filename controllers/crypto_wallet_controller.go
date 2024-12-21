// Package controllers/crypto_wallet_controller.go
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/panaceacode/wallet-demo/models"
	"github.com/panaceacode/wallet-demo/services"
	"net/http"
	"strconv"
	"time"
)

type CryptoWalletController struct {
	walletService         *services.CryptoWalletService
	reconciliationService *services.CryptoReconciliationService
}

func NewCryptoWalletController(walletService *services.CryptoWalletService, reconciliationService *services.CryptoReconciliationService) *CryptoWalletController {
	return &CryptoWalletController{
		walletService:         walletService,
		reconciliationService: reconciliationService,
	}
}

type CryptoCreateWalletRequest struct {
	UserID  uint           `json:"user_id" binding:"required"`
	Network models.Network `json:"network" binding:"required"`
}

type CryptoDepositRequest struct {
	TxHash string `json:"tx_hash" binding:"required"`
}

type CryptoWithdrawRequest struct {
	ToAddress string  `json:"to_address" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}

type CryptoReconciliationRequest struct {
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
}

func (c *CryptoWalletController) CreateWallet(ctx *gin.Context) {
	var req CryptoCreateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := c.walletService.CreateWallet(req.UserID, string(req.Network))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, wallet)
}

func (c *CryptoWalletController) ProcessDeposit(ctx *gin.Context) {
	walletID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet id"})
		return
	}

	var req CryptoDepositRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.walletService.ProcessDeposit(uint(walletID), req.TxHash); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deposit processed successfully"})
}

func (c *CryptoWalletController) Withdraw(ctx *gin.Context) {
	walletID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet id"})
		return
	}

	var req CryptoWithdrawRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txHash, err := c.walletService.Withdraw(uint(walletID), req.ToAddress, req.Amount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "withdrawal initiated successfully",
		"tx_hash": txHash,
	})
}

func (c *CryptoWalletController) GetTransactions(ctx *gin.Context) {
	walletID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet id"})
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	transactions, total, err := c.walletService.GetTransactions(uint(walletID), page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
	})
}

func (c *CryptoWalletController) PerformReconciliation(ctx *gin.Context) {
	walletID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet id"})
		return
	}

	var req CryptoReconciliationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reconciliation, err := c.reconciliationService.PerformReconciliation(
		uint(walletID),
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, reconciliation)
}

func (c *CryptoWalletController) GetReconciliationHistory(ctx *gin.Context) {
	walletID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet id"})
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	reconciliations, total, err := c.reconciliationService.GetReconciliationHistory(
		uint(walletID),
		page,
		pageSize,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reconciliations": reconciliations,
		"total":           total,
		"page":            page,
		"page_size":       pageSize,
	})
}
