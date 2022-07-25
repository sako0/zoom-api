package main

import (
	"log"
	"zoom-api/router"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}
	router := router.Handler()
	// サーバーをポート番号1323で起動
	router.Logger.Fatal(router.Start(":1323"))
}
