package utils

import (
	"user-service/internal/users/models"
)

func StartMigrations() {
	db := GetDBConnection()

	db.AutoMigrate(&models.User{}, &models.Token{})
}
