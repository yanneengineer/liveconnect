package models

import (
	"database/sql"
	"log"
	"os"

	"liveconnect/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = config.Config.DB.Name
	}

	DB, err = sql.Open(config.Config.DB.SQLDriver, dsn)
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
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(8),
		artist TEXT NOT NULL,
		date TIMESTAMP NOT NULL,
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
