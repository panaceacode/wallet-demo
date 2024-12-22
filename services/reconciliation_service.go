package services

import (
	"errors"
	"fmt"
	"github.com/panaceacode/wallet-demo/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type ReconciliationService struct {
	db *gorm.DB
}

func NewReconciliationService(db *gorm.DB) *ReconciliationService {
	return &ReconciliationService{db: db}
}

// PerformReconciliation 执行对账操作
func (s *ReconciliationService) PerformReconciliation(walletID uint, startTime, endTime time.Time, externalBalance decimal.Decimal) (*models.Reconciliation, error) {
	var systemBalance decimal.Decimal

	// 计算系统内的余额变化
	// 1. 获取开始时间之前的最后一个余额
	var lastTx models.Transaction
	err := s.db.Where("wallet_id = ? AND created_at < ?", walletID, startTime).
		Order("created_at DESC").
		First(&lastTx).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	initialBalance := decimal.Zero
	if err == nil {
		initialBalance = lastTx.BalanceAfter
	}

	// 2. 计算时间段内的所有变动
	var transactions []models.Transaction
	err = s.db.Where("wallet_id = ? AND created_at BETWEEN ? AND ?",
		walletID, startTime, endTime).
		Order("created_at ASC").
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	// 3. 计算最终系统余额
	systemBalance = initialBalance
	for _, tx := range transactions {
		if tx.Type == models.TransactionDeposit {
			systemBalance = systemBalance.Add(tx.Amount)
		} else {
			systemBalance = systemBalance.Sub(tx.Amount)
		}
	}

	// 4. 创建对账记录
	difference := systemBalance.Sub(externalBalance)
	status := models.ReconciliationStatusMatched
	if difference.GreaterThan(decimal.NewFromFloat(0.0001)) { // 考虑浮点数精度问题
		status = models.ReconciliationStatusMismatch
	}

	reconciliation := &models.Reconciliation{
		WalletID:        walletID,
		StartTime:       startTime,
		EndTime:         endTime,
		SystemBalance:   systemBalance,
		ExternalBalance: externalBalance,
		Status:          status,
		Difference:      difference,
		Notes:           fmt.Sprintf("Transactions count: %d", len(transactions)),
	}

	err = s.db.Create(reconciliation).Error
	if err != nil {
		return nil, err
	}

	return reconciliation, nil
}

// GetReconciliationHistory 获取对账历史
func (s *ReconciliationService) GetReconciliationHistory(walletID uint, page, pageSize int) ([]models.Reconciliation, error) {
	var reconciliations []models.Reconciliation

	err := s.db.Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reconciliations).Error

	return reconciliations, err
}

// GetReconciliationDetail 获取对账详情
func (s *ReconciliationService) GetReconciliationDetail(reconciliationID uint) (*models.Reconciliation, []models.Transaction, error) {
	var reconciliation models.Reconciliation
	err := s.db.First(&reconciliation, reconciliationID).Error
	if err != nil {
		return nil, nil, err
	}

	var transactions []models.Transaction
	err = s.db.Where("wallet_id = ? AND created_at BETWEEN ? AND ?",
		reconciliation.WalletID, reconciliation.StartTime, reconciliation.EndTime).
		Order("created_at ASC").
		Find(&transactions).Error

	return &reconciliation, transactions, err
}

func (s *WalletService) GetTransactions(walletID uint, startTime, endTime time.Time, page, pageSize int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := s.db.Where("wallet_id = ? AND created_at BETWEEN ? AND ?",
		walletID, startTime, endTime).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&transactions).Error

	return transactions, err
}
