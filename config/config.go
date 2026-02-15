package config

import (
	"log"

	"gopkg.in/ini.v1"
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

// config.ini を読み込む
func LoadConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("config.ini の読み込みに失敗:", err)
	}

	Config = AppConfig{
		Web: WebConfig{
			Port:    cfg.Section("web").Key("port").String(),
			LogFile: cfg.Section("web").Key("logfile").String(),
			Static:  cfg.Section("web").Key("static").String(),
			Views:   cfg.Section("web").Key("views").String(),
		},
		DB: DBConfig{
			SQLDriver: cfg.Section("db").Key("sqldriver").String(),
			Name:      cfg.Section("db").Key("name").String(),
		},
	}
}
