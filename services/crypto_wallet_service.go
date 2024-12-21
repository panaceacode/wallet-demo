// Package services/crypto_wallet_service.go
package services

import (
	"fmt"
	"github.com/panaceacode/wallet-demo/models"
	"gorm.io/gorm"
	"math/big"
)

type CryptoWalletService struct {
	db         *gorm.DB
	blockchain *MockBlockchain
}

func NewCryptoWalletService(db *gorm.DB) *CryptoWalletService {
	return &CryptoWalletService{
		db:         db,
		blockchain: NewMockBlockchain(),
	}
}

// GetBlockchain 返回当前链
func (s *CryptoWalletService) GetBlockchain() *MockBlockchain {
	return s.blockchain
}

// CreateWallet 创建一个钱包
func (s *CryptoWalletService) CreateWallet(userID uint, network string) (*models.CryptoWallet, error) {
	// 创建钱包地址
	address := GenerateAddress()

	var wallet *models.CryptoWallet

	// 开启事务
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 查看在当前链上是不是已经存在一个钱包
		var existingWallet models.CryptoWallet
		err := tx.Where("user_id = ? AND network = ?", userID, network).First(&existingWallet).Error
		if err == nil {
			return fmt.Errorf("wallet already exists for this network")
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// 创建钱包
		wallet = &models.CryptoWallet{
			UserID:    userID,
			Network:   models.Network(network),
			Address:   address,
			Balance:   0,
			ExtraData: "{}", // Initialize empty JSON object
		}

		if err := tx.Create(wallet).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetWallet 获取一个钱包
func (s *CryptoWalletService) GetWallet(walletID uint) (*models.CryptoWallet, error) {
	var wallet models.CryptoWallet
	if err := s.db.First(&wallet, walletID).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetWalletByAddress 根据地址获取钱包
func (s *CryptoWalletService) GetWalletByAddress(address string) (*models.CryptoWallet, error) {
	var wallet models.CryptoWallet
	if err := s.db.Where("address = ?", address).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

// ProcessDeposit 充值
func (s *CryptoWalletService) ProcessDeposit(walletID uint, txHash string) error {
	// 检查这笔交易是否已经处理过了
	var existingTx models.CryptoTransaction
	err := s.db.Where("tx_hash = ?", txHash).First(&existingTx).Error
	if err == nil {
		return fmt.Errorf("transaction already processed")
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	// 从链上获取交易信息
	blockchainTx, err := s.blockchain.GetTransaction(txHash)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %v", err)
	}

	// 校验确认信息
	if blockchainTx.Confirmations < 6 { // 可配置
		return fmt.Errorf("insufficient confirmations: %d/6", blockchainTx.Confirmations)
	}

	// 找到钱包
	var wallet models.CryptoWallet
	if err := s.db.First(&wallet, walletID).Error; err != nil {
		return fmt.Errorf("wallet not found: %v", err)
	}

	// 确认收方地址
	if blockchainTx.To != wallet.Address {
		return fmt.Errorf("invalid recipient address")
	}

	// Convert amount from big.Int to float64
	amount := new(big.Float).SetInt(blockchainTx.Amount)
	finalAmount, _ := amount.Float64()

	// 开启事务落库
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新余额
		if err := tx.Model(&wallet).UpdateColumn(
			"balance",
			gorm.Expr("balance + ?", finalAmount),
		).Error; err != nil {
			return err
		}

		// 创建一条交易信息
		txRecord := &models.CryptoTransaction{
			WalletID:      walletID,
			Type:          models.TransactionDeposit,
			Network:       wallet.Network,
			FromAddress:   blockchainTx.From,
			ToAddress:     wallet.Address,
			Amount:        finalAmount,
			Status:        "completed",
			TxHash:        txHash,
			Confirmations: blockchainTx.Confirmations,
			BlockNumber:   blockchainTx.BlockNumber,
			Raw:           string(blockchainTx.Raw),
		}

		return tx.Create(txRecord).Error
	})
}

// Withdraw 提现功能
func (s *CryptoWalletService) Withdraw(walletID uint, toAddress string, amount float64) (string, error) {
	var (
		txHash string
		err    error
	)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.CryptoWallet
		if err := tx.First(&wallet, walletID).Error; err != nil {
			return fmt.Errorf("wallet not found: %v", err)
		}

		if wallet.Balance < amount {
			return fmt.Errorf("insufficient balance")
		}

		value := new(big.Float).SetFloat64(amount)
		intValue, _ := value.Int(nil)

		hash, err := s.blockchain.SendTransaction(wallet.Address, toAddress, intValue)
		if err != nil {
			return fmt.Errorf("blockchain transaction failed: %v", err)
		}
		txHash = hash

		if err := tx.Model(&wallet).UpdateColumn(
			"balance",
			gorm.Expr("balance - ?", amount),
		).Error; err != nil {
			return fmt.Errorf("failed to update balance: %v", err)
		}

		txRecord := &models.CryptoTransaction{
			WalletID:    walletID,
			Type:        models.TransactionWithdraw,
			Network:     wallet.Network,
			FromAddress: wallet.Address,
			ToAddress:   toAddress,
			Amount:      amount,
			Status:      "processing",
			TxHash:      hash,
		}

		if err := tx.Create(txRecord).Error; err != nil {
			return fmt.Errorf("failed to create transaction record: %v", err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return txHash, nil
}

// GetTransactions 获取交易信息历史
func (s *CryptoWalletService) GetTransactions(walletID uint, page, pageSize int) ([]models.CryptoTransaction, int64, error) {
	var transactions []models.CryptoTransaction
	var total int64

	// 拿到总数
	err := s.db.Model(&models.CryptoTransaction{}).
		Where("wallet_id = ?", walletID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 拿到分页信息
	err = s.db.Where("wallet_id = ?", walletID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}
