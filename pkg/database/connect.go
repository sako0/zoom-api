package database

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open() *gorm.DB {
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database_name := os.Getenv("MYSQL_DATABASE")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database_name + "?charset=utf8mb4&parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("databaseへの接続に失敗しました")
		log.Println(err)
	}
	return db

}
