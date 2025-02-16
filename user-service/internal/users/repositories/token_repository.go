package repositories

import "user-service/internal/users/models"

type TokenRepository interface {
	CreateToken(token *models.Token) error
	DeleteTokenByID(id string) error
	DeleteTokenByUserID(userID string) error
	DeleteTokenByValue(token string) error
	FindTokenByID(id string) (*models.Token, error)
	FindTokenByValue(token string) (*models.Token, error)
	DeleteExpiredTokens() error
}
