package repositories

import (
	"github.com/google/uuid"
	"link-service/internal/links/models"
)

type LinkRepository interface {
	GetAll() ([]*models.Link, error)
	GetAllByUserID(userID uuid.UUID) ([]*models.Link, error)
	GetByID(userID uuid.UUID, id uuid.UUID) (*models.Link, error)
	Save(link *models.Link) error
	GetBySlug(slug string) (*models.Link, error)
	DeleteByID(userID uuid.UUID, id uuid.UUID) error
}
