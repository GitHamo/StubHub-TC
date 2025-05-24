package application

import (
	"github.com/githamo/stubhub-tc/internal/traffic/domain"
	"github.com/google/uuid"
)

type TrafficService struct {
	repo domain.Repository
}

func NewTrafficService(repo domain.Repository) *TrafficService {
	return &TrafficService{
		repo: repo,
	}
}

func (s *TrafficService) GetResponseByUUID(requestId string) (*domain.TrafficResponse, error) {
	if _, err := uuid.Parse(requestId); err != nil {
		return nil, domain.ErrInvalidUUID
	}

	data, err := s.repo.FindByUUID(requestId)

	if err != nil {
		return nil, err
	}

	return domain.NewTrafficResponse(data)
}
