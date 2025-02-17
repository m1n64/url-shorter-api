package utils

import "link-service/internal/links/models"

func StartMigrations() {
	db := GetDBConnection()

	db.AutoMigrate(&models.Link{})
}
