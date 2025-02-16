package repositories

import (
	"gorm.io/gorm"
	"time"
	"user-service/internal/users/models"
)

type tokenGormRepository struct {
	db *gorm.DB
}

func NewTokenGormRepository(db *gorm.DB) TokenRepository {
	return &tokenGormRepository{db: db}
}

func (r *tokenGormRepository) CreateToken(token *models.Token) error {
	return r.db.Create(token).Error
}

func (r *tokenGormRepository) DeleteTokenByID(id string) error {
	return r.db.Delete(&models.Token{}, "id = ?", id).Error
}

func (r *tokenGormRepository) DeleteTokenByUserID(userID string) error {
	return r.db.Delete(&models.Token{}, "user_id = ?", userID).Error
}

func (r *tokenGormRepository) DeleteTokenByValue(token string) error {
	return r.db.Delete(&models.Token{}, "token = ?", token).Error
}

func (r *tokenGormRepository) FindTokenByID(id string) (*models.Token, error) {
	var token models.Token
	if err := r.db.Where("id = ?", id).First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *tokenGormRepository) FindTokenByValue(token string) (*models.Token, error) {
	var tokenModel models.Token
	if err := r.db.Preload("User").Where("token = ?", token).First(&tokenModel).Error; err != nil {
		return nil, err
	}

	return &tokenModel, nil
}

func (r *tokenGormRepository) DeleteExpiredTokens() error {
	return r.db.Delete(&models.Token{}, "expires_at < ?", time.Now()).Error
}
