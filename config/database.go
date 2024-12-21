// Package config config/database.go
package config

import (
	"fmt"
	"github.com/panaceacode/wallet-demo/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
)

type Config struct {
	DBType   string // "mysql" or "sqlite"
	DBPath   string // SQLite database path or MySQL connection string
	DBName   string // Database name for MySQL
	Username string // MySQL username
	Password string // MySQL password
	Host     string // MySQL host
	Port     int    // MySQL port
}

func (c *Config) GetMySQLDSN() string {
	if c.DBPath != "" {
		return c.DBPath
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
}

func NewDB(cfg *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.DBType {
	case "mysql":
		// 构建 root 连接字符串 (不包含数据库名)
		rootDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
		)

		// 先连接到 MySQL 服务器（不指定数据库）
		rootDialector := mysql.Open(rootDSN)
		rootDB, err := gorm.Open(rootDialector, &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MySQL server: %v", err)
		}

		// 创建数据库（如果不存在）
		err = rootDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.DBName)).Error
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %v", err)
		}

		// 构建完整的 DSN（包含数据库名）
		dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
		)

		// 连接到指定的数据库
		dialector = mysql.Open(dbDSN)

	case "sqlite":
		dialector = sqlite.Open(cfg.DBPath)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// Run migrations
	err = migrateDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}

func migrateDB(db *gorm.DB) error {
	// 自动迁移表结构
	baseModels := []interface{}{
		&models.Wallet{},
		&models.Transaction{},
		&models.Reconciliation{},
	}

	// 加密货币钱包系统的表
	cryptoModels := []interface{}{
		&models.CryptoWallet{},
		&models.CryptoTransaction{},
		&models.CryptoReconciliation{},
	}

	// 迁移所有表
	allModels := append(baseModels, cryptoModels...)
	for _, model := range allModels {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %v", model, err)
		}
	}

	// 创建索引
	dialectName := db.Dialector.Name()
	var baseIndexQueries, cryptoIndexQueries []string

	// 基础钱包索引
	switch dialectName {
	case "mysql":
		baseIndexQueries = []string{
			"CREATE INDEX idx_wallets_user_currency ON wallets(user_id, currency)",
			"CREATE INDEX idx_transactions_wallet_created ON transactions(wallet_id, created_at)",
			"CREATE INDEX idx_transactions_tx_hash ON transactions(tx_hash)",
			"CREATE INDEX idx_reconciliations_wallet_time ON reconciliations(wallet_id, start_time, end_time)",
		}
	case "sqlite":
		baseIndexQueries = []string{
			"CREATE INDEX IF NOT EXISTS idx_wallets_user_currency ON wallets(user_id, currency)",
			"CREATE INDEX IF NOT EXISTS idx_transactions_wallet_created ON transactions(wallet_id, created_at)",
			"CREATE INDEX IF NOT EXISTS idx_transactions_tx_hash ON transactions(tx_hash)",
			"CREATE INDEX IF NOT EXISTS idx_reconciliations_wallet_time ON reconciliations(wallet_id, start_time, end_time)",
		}
	}

	// 加密货币钱包索引
	switch dialectName {
	case "mysql":
		cryptoIndexQueries = []string{
			"CREATE INDEX idx_crypto_wallets_user_network ON crypto_wallets(user_id, network)",
			"CREATE INDEX idx_crypto_wallets_address ON crypto_wallets(address)",
			"CREATE INDEX idx_crypto_transactions_wallet ON crypto_transactions(wallet_id, created_at)",
			"CREATE INDEX idx_crypto_transactions_addresses ON crypto_transactions(from_address, to_address)",
			"CREATE INDEX idx_crypto_reconciliations_wallet ON crypto_reconciliations(wallet_id, start_time, end_time)",
		}
	case "sqlite":
		cryptoIndexQueries = []string{
			"CREATE INDEX IF NOT EXISTS idx_crypto_wallets_user_network ON crypto_wallets(user_id, network)",
			"CREATE INDEX IF NOT EXISTS idx_crypto_wallets_address ON crypto_wallets(address)",
			"CREATE INDEX IF NOT EXISTS idx_crypto_transactions_wallet ON crypto_transactions(wallet_id, created_at)",
			"CREATE INDEX IF NOT EXISTS idx_crypto_transactions_addresses ON crypto_transactions(from_address, to_address)",
			"CREATE INDEX IF NOT EXISTS idx_crypto_reconciliations_wallet ON crypto_reconciliations(wallet_id, start_time, end_time)",
		}
	}

	// 执行所有索引创建
	allQueries := append(baseIndexQueries, cryptoIndexQueries...)
	for _, query := range allQueries {
		if err := db.Exec(query).Error; err != nil {
			if dialectName == "mysql" && strings.Contains(err.Error(), "Duplicate key name") {
				continue
			}
			return fmt.Errorf("failed to create index: %v", err)
		}
	}

	return nil
}
