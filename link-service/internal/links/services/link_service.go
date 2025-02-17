package services

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"link-service/internal/links/models"
	"link-service/internal/links/repositories"
)

type LinkService struct {
	linkRepository   repositories.LinkRepository
	linkCacheService *LinkCacheService
	db               *gorm.DB
	logger           *zap.Logger
}

func NewLinkService(linkRepository repositories.LinkRepository, linkCacheService *LinkCacheService, db *gorm.DB, logger *zap.Logger) *LinkService {
	return &LinkService{
		linkRepository:   linkRepository,
		linkCacheService: linkCacheService,
		db:               db,
		logger:           logger,
	}
}

func (s *LinkService) GetAll() ([]*models.Link, error) {
	return s.linkRepository.GetAll()
}

func (s *LinkService) GetAllByUserID(userID uuid.UUID) ([]*models.Link, error) {
	return s.linkRepository.GetAllByUserID(userID)
}

func (s *LinkService) Create(userID uuid.UUID, link string, slug string) (*models.Link, error) {
	linkModel := &models.Link{
		UserID: userID,
		Slug:   slug,
		URL:    link,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		err := s.linkRepository.Save(linkModel)
		if err != nil {
			return err
		}

		err = s.linkCacheService.SaveLinkInGlobalCache(linkModel)
		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
	if err != nil {
		s.logger.Error("Error creating link: ", zap.Error(err))
		return nil, err
	}

	return linkModel, nil
}
