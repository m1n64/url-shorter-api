package repositories

import (
	"user-service/internal/users/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(email string) (*models.User, error)
	FindByToken(token string) (*models.User, error)
}
