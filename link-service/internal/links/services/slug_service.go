package services

import "link-service/pkg/utils"

const MaxSlugLength = 6

type SlugService struct {
}

func NewSlugService() *SlugService {
	return &SlugService{}
}

func (s *SlugService) GenerateSlug() string {
	return utils.RandStringBytesRmndr(MaxSlugLength)
}
