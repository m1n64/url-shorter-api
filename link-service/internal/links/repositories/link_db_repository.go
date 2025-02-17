package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"link-service/internal/links/models"
)

type linkDBRepository struct {
	db *gorm.DB
}

func NewLinkDBRepository(db *gorm.DB) LinkRepository {
	return &linkDBRepository{
		db: db,
	}
}

func (r *linkDBRepository) GetAll() ([]*models.Link, error) {
	var links []*models.Link
	if err := r.db.Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (r *linkDBRepository) GetAllByUserID(userID uuid.UUID) ([]*models.Link, error) {
	var links []*models.Link
	if err := r.db.Where("user_id = ?", userID).Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (r *linkDBRepository) GetByID(userID uuid.UUID, id uuid.UUID) (*models.Link, error) {
	var link models.Link
	if err := r.db.Where("user_id = ? AND id = ?", userID, id).First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *linkDBRepository) Save(link *models.Link) error {
	return r.db.Save(link).Error
}

func (r *linkDBRepository) GetBySlug(slug string) (*models.Link, error) {
	var link models.Link
	if err := r.db.Where("slug = ?", slug).First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *linkDBRepository) DeleteByID(userID uuid.UUID, id uuid.UUID) error {
	return r.db.Where("user_id = ? AND id = ?", userID, id).Delete(&models.Link{}).Error
}
