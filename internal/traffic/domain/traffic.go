package domain

import (
	"encoding/json"
	"errors"
)

var (
	ErrTrafficDataNotFound = errors.New("traffic data not found")
	ErrInvalidUUID         = errors.New("invalid UUID format")
	ErrInvalidContent      = errors.New("invalid JSON content")
)

type TrafficResponseContent []byte

type TrafficResponse struct {
	Content json.RawMessage
}

func NewTrafficResponse(content json.RawMessage) (*TrafficResponse, error) {
	if !json.Valid(content) {
		return nil, ErrInvalidContent
	}

	return &TrafficResponse{
		Content: []byte(content),
	}, nil
}

type Service interface {
	GetResponseByUUID(uuid string) (*TrafficResponse, error)
}

type Repository interface {
	FindByUUID(uuid string) (json.RawMessage, error)
}
