package main

import (
	"liveconnect/app/controllers"
	"liveconnect/app/models"
	"liveconnect/config"
)

func main() {
	// 設定読み込み
	config.LoadConfig()

	// DB初期化
	models.InitDB()

	// サーバ起動
	if err := controllers.StartMainServer(); err != nil {
		panic(err)
	}
}
