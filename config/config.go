package config

import (
	"os"
)

// Web設定
type WebConfig struct {
	Port    string
	LogFile string
	Static  string
	Views   string
}

// DB設定
type DBConfig struct {
	SQLDriver string
	Name      string
}

// アプリ全体の設定
type AppConfig struct {
	Web WebConfig
	DB  DBConfig
}

// グローバルで参照する設定
var Config AppConfig

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadConfig は環境変数から設定を読み込む
func LoadConfig() {

	Config = AppConfig{
		Web: WebConfig{
			Port:    getEnv("WEB_PORT", "8080"),
			LogFile: getEnv("WEB_LOGFILE", "webapp.log"),
			Static:  getEnv("WEB_STATIC", "app/static"),
			Views:   getEnv("WEB_VIEWS", "app/views"),
		},
		DB: DBConfig{
			SQLDriver: getEnv("DB_DRIVER", "postgres"),
			Name:      getEnv("DB_NAME", ""),
		},
	}
}
