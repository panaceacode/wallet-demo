package main

import (
	"github.com/gin-gonic/gin"
	"github.com/panaceacode/wallet-demo/db"
	"github.com/panaceacode/wallet-demo/handlers"
)

func main() {
	// 初始化数据库
	db.InitDB()
	defer db.CloseDB()

	// 启动 Gin Web 服务
	r := gin.Default()

	// 注册路由
	r.GET("/balance", handlers.GetBalanceHandler)
	r.POST("/deposit", handlers.DepositHandler)
	r.POST("/withdraw", handlers.WithdrawHandler)

	// 启动服务
	r.Run(":8080")
}
