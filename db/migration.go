package db

// 创建数据库表
func CreateTables() {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS wallet (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			currency TEXT NOT NULL,
			balance REAL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_user_currency UNIQUE (user_id, currency)
		);
	`)
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS utxo (
			tx_id TEXT,
			output_index INTEGER,
			user_id INTEGER,
			amount REAL,
			is_spent BOOLEAN,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY(tx_id, output_index)
		);
	`)
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS transaction_record (
			tx_id TEXT PRIMARY KEY,
			user_id INTEGER,
			currency TEXT,
			amount REAL,
			status TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		panic(err)
	}
}
