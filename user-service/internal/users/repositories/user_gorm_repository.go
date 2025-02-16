package repositories

import (
	"gorm.io/gorm"
	"user-service/internal/users/models"
)

type userGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) UserRepository {
	return &userGormRepository{db: db}
}

func (r *userGormRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *userGormRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userGormRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userGormRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userGormRepository) FindByToken(token string) (*models.User, error) {
	var user models.User
	if err := r.db.Joins("JOIN tokens ON tokens.user_id = users.id").
		Where("tokens.token = ?", token).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
