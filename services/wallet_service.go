// Package services/wallet_service.go
package services

import (
	"errors"
	"github.com/panaceacode/wallet-demo/models"
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
		Balance:  0,
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

func (s *WalletService) Deposit(walletID uint, amount float64, txHash string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, walletID).Error; err != nil {
			return err
		}

		balanceBefore := wallet.Balance
		wallet.Balance += amount

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

func (s *WalletService) Withdraw(walletID uint, amount float64, txHash string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, walletID).Error; err != nil {
			return err
		}

		if wallet.Balance < amount {
			return errors.New("insufficient balance")
		}

		balanceBefore := wallet.Balance
		wallet.Balance -= amount

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
