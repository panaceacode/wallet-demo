package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// 初始化数据库连接
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", ":memory:") // 使用内存数据库
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// 创建表结构
	CreateTables()
}

// 关闭数据库连接
func CloseDB() {
	DB.Close()
}
