package services

import (
	"encoding/json"
	"fmt"
	"github.com/panaceacode/wallet-demo/models"
	"gorm.io/gorm"
	"math"
	"math/big"
	"strings"
	"time"
)

type CryptoReconciliationService struct {
	db         *gorm.DB
	blockchain *MockBlockchain
}

func NewCryptoReconciliationService(db *gorm.DB, blockchain *MockBlockchain) *CryptoReconciliationService {
	return &CryptoReconciliationService{
		db:         db,
		blockchain: blockchain,
	}
}

// PerformReconciliation 执行链上数据对账
func (s *CryptoReconciliationService) PerformReconciliation(walletID uint, startTime, endTime time.Time) (*models.CryptoReconciliation, error) {
	var wallet models.CryptoWallet
	if err := s.db.First(&wallet, walletID).Error; err != nil {
		return nil, fmt.Errorf("wallet not found: %v", err)
	}

	// 获取系统中记录的交易
	var systemTransactions []models.CryptoTransaction
	err := s.db.Where("wallet_id = ? AND created_at BETWEEN ? AND ?",
		walletID, startTime, endTime).
		Find(&systemTransactions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get system transactions: %v", err)
	}

	// 获取链上交易记录
	chainTransactions, err := s.blockchain.GetTransactionHistory(wallet.Address, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get blockchain transactions: %v", err)
	}

	// 获取链上余额
	chainBalance, err := s.blockchain.GetAddressBalance(wallet.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get blockchain balance: %v", err)
	}
	chainBalanceFloat := new(big.Float).SetInt(chainBalance)
	finalChainBalance, _ := chainBalanceFloat.Float64()

	// 创建对账记录
	reconciliation := &models.CryptoReconciliation{
		WalletID:      walletID,
		StartTime:     startTime,
		EndTime:       endTime,
		SystemBalance: wallet.Balance,
		ChainBalance:  finalChainBalance,
		Status:        models.ReconciliationStatusMatched,
		Difference:    wallet.Balance - finalChainBalance,
	}

	// 分析差异
	if math.Abs(reconciliation.Difference) > 0.0001 {
		reconciliation.Status = models.ReconciliationStatusMismatch
		s.analyzeMismatch(reconciliation, systemTransactions, chainTransactions)
	}

	// 保存对账记录
	if err := s.db.Create(reconciliation).Error; err != nil {
		return nil, err
	}

	return reconciliation, nil
}

func (s *CryptoReconciliationService) analyzeMismatch(
	reconciliation *models.CryptoReconciliation,
	systemTxs []models.CryptoTransaction,
	chainTxs []*BlockchainTransaction) {

	// 构建交易映射
	chainTxMap := make(map[string]*BlockchainTransaction)
	for _, tx := range chainTxs {
		chainTxMap[tx.Hash] = tx
	}

	// 查找未匹配的交易
	var unmatchedTxs []string
	var reasons []string

	for _, sysTx := range systemTxs {
		chainTx, exists := chainTxMap[sysTx.TxHash]
		if !exists {
			unmatchedTxs = append(unmatchedTxs, sysTx.TxHash)
			reasons = append(reasons, fmt.Sprintf("Transaction %s not found on chain", sysTx.TxHash))
			continue
		}

		// 比较金额
		sysAmount := new(big.Float).SetFloat64(sysTx.Amount)
		chainAmount := new(big.Float).SetInt(chainTx.Amount)
		if sysAmount.Cmp(chainAmount) != 0 {
			reasons = append(reasons, fmt.Sprintf(
				"Amount mismatch for tx %s: system=%v, chain=%v",
				sysTx.TxHash,
				sysTx.Amount,
				chainTx.Amount,
			))
		}
	}

	if len(unmatchedTxs) > 0 {
		unmatchedJSON, _ := json.Marshal(unmatchedTxs)
		reconciliation.UnmatchedTxs = string(unmatchedJSON)
	}
	if len(reasons) > 0 {
		reconciliation.MismatchReason = strings.Join(reasons, "; ")
	}
}

// GetReconciliationHistory 获取对账历史
func (s *CryptoReconciliationService) GetReconciliationHistory(walletID uint, page, pageSize int) ([]models.CryptoReconciliation, int64, error) {
	var reconciliations []models.CryptoReconciliation
	var total int64

	err := s.db.Model(&models.CryptoReconciliation{}).
		Where("wallet_id = ?", walletID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reconciliations).Error

	return reconciliations, total, err
}
