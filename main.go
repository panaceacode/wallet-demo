package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/panaceacode/wallet-demo/config"
	"github.com/panaceacode/wallet-demo/controllers"
	"github.com/panaceacode/wallet-demo/services"
)

func main() {
	cfg := &config.Config{
		DBType:   "mysql",
		DBName:   "wallet",
		Username: "root",   // 替换为你的 MySQL 用户名
		Password: "970827", // 替换为你的 MySQL 密码
		Host:     "localhost",
		Port:     3306,
	}
	// 或者使用 SQLite
	// cfg := &config.Config{
	//     DBType: "sqlite",
	//     DBPath: "wallet.db",
	// }

	db, err := config.NewDB(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize database: %v", err))
	}

	// Initialize services
	walletService := services.NewWalletService(db)
	reconciliationService := services.NewReconciliationService(db)
	walletController := controllers.NewWalletController(walletService, reconciliationService)

	cryptoWalletService := services.NewCryptoWalletService(db)
	cryptoReconciliationService := services.NewCryptoReconciliationService(db, cryptoWalletService.GetBlockchain())
	cryptoWalletController := controllers.NewCryptoWalletController(cryptoWalletService, cryptoReconciliationService)

	r := gin.Default()

	// Routes
	api := r.Group("/api")
	{
		wallets := api.Group("/wallets")
		{
			wallets.POST("/", walletController.CreateWallet)
			wallets.POST("/:id/deposit", walletController.Deposit)
			wallets.POST("/:id/withdraw", walletController.Withdraw)
			wallets.GET("/:id/transactions", walletController.GetTransactions)

			wallets.POST("/:id/reconciliation", walletController.PerformReconciliation)
			wallets.GET("/:id/reconciliation/history", walletController.GetReconciliationHistory)
			wallets.GET("/reconciliation/:id", walletController.GetReconciliationDetail)
		}

		// 加密货币钱包路由
		cryptoWallets := api.Group("/crypto-wallets")
		{
			cryptoWallets.POST("/", cryptoWalletController.CreateWallet)
			cryptoWallets.POST("/:id/deposit", cryptoWalletController.ProcessDeposit)
			cryptoWallets.POST("/:id/withdraw", cryptoWalletController.Withdraw)
			cryptoWallets.GET("/:id/transactions", cryptoWalletController.GetTransactions)
			cryptoWallets.POST("/:id/reconciliation", cryptoWalletController.PerformReconciliation)
			cryptoWallets.GET("/:id/reconciliation/history", cryptoWalletController.GetReconciliationHistory)
		}

	}

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
