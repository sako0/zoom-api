package main

import (
	"zoom-api/pkg/database"
	"zoom-api/pkg/models"
)

func main() {
	db := database.Open()

	db.AutoMigrate(
		&models.User{},
		&models.Organization{},
		&models.OrganizationMember{},
	)
	organization := []models.Organization{{Name: "aaa"}, {Name: "bbb"}, {Name: "ccc"}}

	db.Create(&organization)
}
