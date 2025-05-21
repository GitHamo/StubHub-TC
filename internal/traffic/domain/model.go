package domain

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrTrafficDataNotFound = errors.New("traffic data not found")
	ErrInvalidUUID         = errors.New("invalid UUID format")
	ErrInvalidContent      = errors.New("invalid JSON content")
)

type TrafficEndpointData struct {
	Path string
}

type TrafficContentData struct {
	Filename string
	Content  json.RawMessage
}

func NewTrafficData(trafficId string, content json.RawMessage) (*TrafficContentData, error) {
	if _, err := uuid.Parse(trafficId); err != nil {
		return nil, ErrInvalidUUID
	}

	if !json.Valid(content) {
		return nil, ErrInvalidContent
	}

	return &TrafficContentData{
		Filename: trafficId,
		Content:  content,
	}, nil
}

type Repository interface {
	FindByUUID(uuid string) (*TrafficContentData, error)
}
