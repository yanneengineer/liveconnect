package models

import (
	"database/sql"
	"log"

	"liveconnect/config"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open(
		config.Config.DB.SQLDriver,
		config.Config.DB.Name,
	)
	if err != nil {
		log.Fatal("DB接続失敗:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB Ping 失敗:", err)
	}

	createTables()
}

func createTables() {
	// lives テーブル
	createLivesTable := `
	CREATE TABLE IF NOT EXISTS lives (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id VARCHAR(8),
		artist TEXT NOT NULL,
		date TEXT NOT NULL,
		venue TEXT NOT NULL
	);
	`

	// users テーブル
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(8) PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);
	`

	if _, err := DB.Exec(createLivesTable); err != nil {
		log.Fatal("livesテーブル作成失敗:", err)
	}

	if _, err := DB.Exec(createUsersTable); err != nil {
		log.Fatal("usersテーブル作成失敗:", err)
	}
}
