// Package services/wallet_service.go
package services

import (
	"errors"
	"github.com/panaceacode/wallet-demo/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletService struct {
	db *gorm.DB
}

func NewWalletService(db *gorm.DB) *WalletService {
	return &WalletService{db: db}
}

func (s *WalletService) CreateWallet(userID uint, currency string) (*models.Wallet, error) {
	wallet := &models.Wallet{
		UserID:   userID,
		Currency: currency,
		Balance:  decimal.NewFromFloat(0),
	}

	err := s.db.Create(wallet).Error
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletService) GetWallet(userID uint, currency string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := s.db.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (s *WalletService) Deposit(walletID uint, amount decimal.Decimal, txHash string) error {
	if amount.LessThan(decimal.Zero) {
		return errors.New("invalid deposit amount")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, walletID).Error; err != nil {
			return err
		}

		// 检查交易哈希是否已存在
		var existingTx models.Transaction
		if err := tx.Where("tx_hash = ?", txHash).First(&existingTx).Error; err == nil {
			return errors.New("transaction already processed")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		balanceBefore := wallet.Balance
		wallet.Balance = wallet.Balance.Add(amount)

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		transaction := models.Transaction{
			WalletID:      walletID,
			Type:          models.TransactionDeposit,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  wallet.Balance,
			Status:        "completed",
			TxHash:        txHash,
		}

		return tx.Create(&transaction).Error
	})
}

func (s *WalletService) Withdraw(walletID uint, amount decimal.Decimal, txHash string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, walletID).Error; err != nil {
			return err
		}

		if wallet.Balance.LessThan(amount) {
			return errors.New("insufficient balance")
		}

		// 检查交易哈希是否已存在
		var existingTx models.Transaction
		if err := tx.Where("tx_hash = ?", txHash).First(&existingTx).Error; err == nil {
			return errors.New("transaction already processed")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		balanceBefore := wallet.Balance
		wallet.Balance = wallet.Balance.Sub(amount)

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		transaction := models.Transaction{
			WalletID:      walletID,
			Type:          models.TransactionWithdraw,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  wallet.Balance,
			Status:        "completed",
			TxHash:        txHash,
		}

		return tx.Create(&transaction).Error
	})
}
