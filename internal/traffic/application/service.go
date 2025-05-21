package application

import (
	"encoding/json"

	"github.com/githamo/stubhub-tc/internal/traffic/domain"
)

type TrafficService struct {
	repo domain.Repository
}

func NewTrafficService(repo domain.Repository) *TrafficService {
	return &TrafficService{
		repo: repo,
	}
}

func (s *TrafficService) GetContentByUUID(uuid string) (json.RawMessage, error) {
	data, err := s.repo.FindByUUID(uuid)

	if err != nil {
		return nil, err
	}

	return data.Content, nil
}
