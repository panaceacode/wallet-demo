package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"
)

// BlockchainTransaction represents a transaction on the blockchain
type BlockchainTransaction struct {
	Hash          string
	From          string
	To            string
	Amount        *big.Int
	BlockNumber   uint64
	Confirmations int
	Timestamp     time.Time
	Status        string // success, pending, failed
	Fee           *big.Int
	Raw           []byte // Raw transaction data
}

// MockBlockchain simulates a blockchain network
type MockBlockchain struct {
	mutex        sync.RWMutex
	transactions map[string]*BlockchainTransaction
	balances     map[string]*big.Int
	currentBlock uint64
}

func NewMockBlockchain() *MockBlockchain {
	return &MockBlockchain{
		transactions: make(map[string]*BlockchainTransaction),
		balances:     make(map[string]*big.Int),
		currentBlock: 0,
	}
}

// GetTransaction retrieves a specific transaction by hash
func (b *MockBlockchain) GetTransaction(txHash string) (*BlockchainTransaction, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	tx, exists := b.transactions[txHash]
	if !exists {
		return nil, fmt.Errorf("transaction not found")
	}
	return tx, nil
}

// SendTransaction simulates sending a transaction to the blockchain
func (b *MockBlockchain) SendTransaction(from, to string, amount *big.Int) (string, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check sender balance
	balance, exists := b.balances[from]
	if !exists {
		balance = big.NewInt(0)
	}

	if balance.Cmp(amount) < 0 {
		return "", fmt.Errorf("insufficient balance")
	}

	// Generate transaction hash
	txHash := generateTxHash()

	// Create transaction record
	tx := &BlockchainTransaction{
		Hash:          txHash,
		From:          from,
		To:            to,
		Amount:        amount,
		BlockNumber:   b.currentBlock + 1,
		Confirmations: 0,
		Timestamp:     time.Now(),
		Status:        "pending",
		Fee:           big.NewInt(21000), // Mock gas fee
		Raw:           []byte(fmt.Sprintf("mock_tx_data_%s", txHash)),
	}

	// Update balances
	b.balances[from] = new(big.Int).Sub(balance, amount)
	toBalance, exists := b.balances[to]
	if !exists {
		toBalance = big.NewInt(0)
	}
	b.balances[to] = new(big.Int).Add(toBalance, amount)

	// Store transaction
	b.transactions[txHash] = tx
	b.currentBlock++

	// Start confirmation simulation
	go b.simulateConfirmations(txHash)

	return txHash, nil
}

func (b *MockBlockchain) simulateConfirmations(txHash string) {
	for i := 1; i <= 12; i++ {
		time.Sleep(5 * time.Second) // Simulate block time

		b.mutex.Lock()
		if tx, exists := b.transactions[txHash]; exists {
			tx.Confirmations = i
			if i >= 6 { // After 6 confirmations
				tx.Status = "success"
			}
		}
		b.mutex.Unlock()
	}
}

// GetTransactionHistory returns all transactions for an address within a time range
func (b *MockBlockchain) GetTransactionHistory(address string, startTime, endTime time.Time) ([]*BlockchainTransaction, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	var transactions []*BlockchainTransaction
	for _, tx := range b.transactions {
		if (tx.From == address || tx.To == address) &&
			tx.Timestamp.After(startTime) &&
			tx.Timestamp.Before(endTime) {
			transactions = append(transactions, tx)
		}
	}
	return transactions, nil
}

// GetAddressBalance returns the current balance of an address from blockchain
func (b *MockBlockchain) GetAddressBalance(address string) (*big.Int, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	balance, exists := b.balances[address]
	if !exists {
		return big.NewInt(0), nil
	}
	return balance, nil
}

// 生成模拟的交易哈希
func generateTxHash() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return "0x" + hex.EncodeToString(bytes)
}

// 生成模拟的钱包地址
func GenerateAddress() string {
	bytes := make([]byte, 20)
	rand.Read(bytes)
	return "0x" + hex.EncodeToString(bytes)
}
